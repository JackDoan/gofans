package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gofans/hwmonHelp"
)

type FanSetting struct {
	Name  string
	Value uint
}

type FanCurvePoint struct {
	Temp uint
	Pwm  uint8
}

func (s *FanSetting) Set(devicePath string) error {
	settingPath := filepath.Join(devicePath, s.Name)
	f, err := os.OpenFile(settingPath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", settingPath, err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()
	out := fmt.Sprintf("%d\n", s.Value)
	log.Printf("setting %s to %s\n", settingPath, out)
	_, err = f.WriteString(out)
	if err != nil {
		return fmt.Errorf("could not write to %s: %w", settingPath, err)
	}
	return err
}

type FanCurve struct {
	PwmNumber int
	Points    [5]FanCurvePoint
}

func (s *FanCurve) Set(devicePath string) error {
	var errs []error
	for i, v := range s.Points {
		t := FanSetting{fmt.Sprintf("pwm%d_auto_point%d_temp", s.PwmNumber, i+1), v.Temp}
		p := FanSetting{fmt.Sprintf("pwm%d_auto_point%d_pwm", s.PwmNumber, i+1), uint(v.Pwm)}
		err := t.Set(devicePath)
		if err != nil {
			errs = append(errs, err)
		}
		err = p.Set(devicePath)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func setMotherboard() error {
	const motherboardFancontrol = "nct6799"
	nct, err := hwmonHelp.FindDeviceByName(motherboardFancontrol) //nct6799-isa-0290
	if err != nil {
		return fmt.Errorf("could not find motherboard fan controller %s: %w", motherboardFancontrol, err)
	}
	vrmFanPoints := [5]FanCurvePoint{
		{45000, 0},
		{55000, 102},
		{70000, 178},
		{74000, 252},
		{75000, 255},
	}
	var errs []error
	curves := []FanCurve{
		{
			PwmNumber: 5, //front fans?
			Points: [5]FanCurvePoint{
				{20000, 150},
				{70000, 178},
				{75000, 255},
				{75000, 255},
				{100000, 255},
			},
		},
		{PwmNumber: 6 /* VRM1 */, Points: vrmFanPoints},
		{PwmNumber: 7 /* VRM2 */, Points: vrmFanPoints},
	}
	for _, s := range curves {
		err = s.Set(nct)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func setPump() error {
	const pumpName = "d5next"

	p, err := hwmonHelp.FindDeviceByName(pumpName)
	if err != nil {
		return fmt.Errorf("could not find pump controller %s: %w", pumpName, err)
	}

	var errs []error
	settings := []FanSetting{
		{"pwm1", 85},  //pump
		{"pwm2", 186}, //rad fans
	}
	for _, s := range settings {
		err = s.Set(p)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func main() {
	exitCode := 0
	err := setMotherboard()
	if err != nil {
		exitCode |= 1
		log.Printf("could not set motherboard fans: %v", err)
	}
	err = setPump()
	if err != nil {
		exitCode |= 1 << 1
		log.Printf("could not set pump fans: %v", err)
	}
	os.Exit(exitCode)
}
