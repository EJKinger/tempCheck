package main

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
)

var (
	deviceDirPath        = kingpin.Flag("devicePath", "Path to directory with device info.").Short('d').Required().ExistingDir()
	textfileExporterPath = kingpin.Flag("textfileExporterPath", "Path to directory for node_exporter textfile collector.").Short('t').Required().File()
)

type ProbeTemp struct {
	C     float64
	Probe string
}

func getTemps(d string) []*ProbeTemp {
	temps := []*ProbeTemp{}
	files, err := ioutil.ReadDir(d)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		dir := path.Base(file.Name())
		if dir != "w1_bus_master1" {
			content, err := ioutil.ReadFile(path.Join(d, file.Name(), "w1_slave"))
			if err != nil {
				fmt.Println(err)
			} else {
				c, err := getTempC(content)
				if err != nil {
					fmt.Println(err)
					continue
				}
				pt := ProbeTemp{
					Probe: dir,
					C:     c,
				}
				temps = append(temps, &pt)
			}
		}

	}
	return temps

}

func writeExporterStrings(f *os.File, ts []*ProbeTemp) error {
	tf, err := ioutil.TempFile(filepath.Dir(f.Name()), filepath.Base(f.Name()))
	if err != nil {
		return err
	}
	metrics := ""
	for _, v := range ts {
		metrics += fmt.Sprintf("probe_temp_celsius{probe=\"%s\",} %f\n", v.Probe, v.C)
	}
	_, err = io.WriteString(tf, metrics)
	tf.Close()
	if err != nil {
		return err
	}
	err = os.Rename(tf.Name(), f.Name())
	if err != nil {
		return err
	}
	return nil
}

func getTempC(c []byte) (float64, error) {
	r, _ := regexp.Compile("t=(\\d*)")
	m := r.FindStringSubmatch(string(c))
	if m == nil || len(m) < 2 {
		return 0, errors.New("bad sensor read")
	}
	v, err := strconv.ParseFloat(string(m[1]), 64)
	if err != nil {
		return 0, errors.New("bad sensor read")
	}

	return v / 1000.0, nil
}

func main() {
	kingpin.Version("0.1.0")
	kingpin.Parse()
	ts := getTemps(*deviceDirPath)
	err := writeExporterStrings(*textfileExporterPath, ts)
	if err != nil {
		fmt.Println(err)
	}
}
