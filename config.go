/*
 * Copyright (c) 2018 Manabu Sonoda
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dtap

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

type Config struct {
	InputMsgBuffer uint
	InputUnix      []*InputUnixSocketConfig
	InputFile      []*InputFileConfig
	InputTail      []*InputTailConfig
	InputTCP       []*InputTCPSocketConfig
	OutputUnix     []*OutputUnixSocketConfig
	OutputFile     []*OutputFileConfig
	OutputTCP      []*OutputTCPSocketConfig
	OutputFluent   []*OutputFluentConfig
}

func (c *Config) Validate() []error {
	errs := []error{}
	if c.InputMsgBuffer < 128 {
		errs = append(errs, errors.New("InputMsgBuffer must not small 128"))
	}
	for n, i := range c.InputUnix {
		if err := i.Validate(); err != nil {
			err.configType = "InputUnix"
			err.no = n
			errs = append(errs, err)
		}
	}
	for n, i := range c.InputFile {
		if err := i.Validate(); err != nil {
			err.configType = "InputFile"
			err.no = n
			errs = append(errs, err)
		}
	}
	for n, i := range c.InputTCP {
		if err := i.Validate(); err != nil {
			err.configType = "InputTCP"
			err.no = n
			errs = append(errs, err)
		}
	}
	for n, o := range c.OutputUnix {
		if err := o.Validate(); err != nil {
			err.configType = "OutputUnix"
			err.no = n
			errs = append(errs, err)
		}
	}
	for n, o := range c.OutputFile {
		if err := o.Validate(); err != nil {
			err.configType = "OutputFile"
			err.no = n
			errs = append(errs, err)
		}
	}
	for n, o := range c.OutputTCP {
		if err := o.Validate(); err != nil {
			err.configType = "OutputTCP"
			err.no = n
			errs = append(errs, err)
		}
	}
	for n, o := range c.OutputFluent {
		if err := o.Validate(); err != nil {
			err.configType = "OutputFluent"
			err.no = n
			errs = append(errs, err)
		}
	}
	return errs
}

type ValidationError struct {
	configType string
	no         int
	errors     []error
}

func (e *ValidationError) Error() string {
	var msg string
	for _, err := range e.errors {
		msg += fmt.Sprintf("%s[%d]: %s\n", e.configType, e.no, err.Error())
	}
	return msg
}

func (e *ValidationError) Add(err error) {
	e.errors = append(e.errors, err)
}

func (e *ValidationError) Err() *ValidationError {
	if len(e.errors) > 0 {
		return e
	}
	return nil
}

func NewValidationError() *ValidationError {
	return &ValidationError{
		errors: []error{},
	}
}

func NewConfigFromFile(filename string) (*Config, error) {
	c := &Config{}
	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigType("toml")
	v.SetDefault("InputMsgBuffer", 10000)
	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "can't read config")
	}
	if err := v.Unmarshal(c); err != nil {
		return nil, errors.Wrap(err, "can't parse config")
	}
	return c, nil
}

type InputUnixSocketConfig struct {
	Path string
	User string
}

func (i *InputUnixSocketConfig) Validate() *ValidationError {
	err := NewValidationError()
	if i.Path == "" {
		err.Add(errors.New("Path must not be empty"))
	}
	return err.Err()
}

func (i *InputUnixSocketConfig) GetPath() string {
	return i.Path
}
func (i *InputUnixSocketConfig) GetUser() string {
	return i.User
}

type InputFileConfig struct {
	Path string
}

func (i *InputFileConfig) Validate() *ValidationError {
	err := NewValidationError()
	if i.Path == "" {
		err.Add(errors.New("Path must not be empty"))
	}
	return err.Err()
}

func (i *InputFileConfig) GetPath() string {
	return i.Path
}

type InputTailConfig struct {
	Path string
}

func (i *InputTailConfig) Validate() *ValidationError {
	err := NewValidationError()
	if i.Path == "" {
		err.Add(errors.New("Path must not be empty"))
	}
	return err.Err()
}

func (i *InputTailConfig) GetPath() string {
	return i.Path
}

type InputTCPSocketConfig struct {
	Address string
	Port    uint16
}

func (i *InputTCPSocketConfig) Validate() *ValidationError {
	err := NewValidationError()
	if i.Address == "" {
		err.Add(errors.New("Host must not be empty"))
	}
	return err.Err()
}

func (i *InputTCPSocketConfig) GetNet() string {
	address := i.Address
	port := i.Port
	if address == "" {
		address = "0.0.0.0"
	}
	if port == 0 {
		port = 10053
	}
	if strings.Contains(address, ":") {
		address = "[" + address + "]"
	}
	return address + ":" + strconv.Itoa(int(port))
}

type OutputUnixSocketConfig struct {
	Path       string
	BufferSize uint
}

func (o *OutputUnixSocketConfig) Validate() *ValidationError {
	err := NewValidationError()
	if o.Path == "" {
		err.Add(errors.New("Path must not be empty"))
	}
	return err.Err()
}

func (o *OutputUnixSocketConfig) GetPath() string {
	return o.Path
}

func (o *OutputUnixSocketConfig) GetBufferSize() uint {
	if o.BufferSize == 0 {
		return OutputBufferSize
	}
	return o.BufferSize
}

type OutputFileConfig struct {
	Path       string
	User       string
	BufferSize uint
}

func (o *OutputFileConfig) Validate() *ValidationError {
	err := NewValidationError()
	if o.Path == "" {
		err.Add(errors.New("Path must not be empty"))
	}
	return err.Err()
}

func (o *OutputFileConfig) GetPath() string {
	return o.Path
}
func (o *OutputFileConfig) GetUser() string {
	return o.User
}

func (o *OutputFileConfig) GetBufferSize() uint {
	if o.BufferSize == 0 {
		return OutputBufferSize
	}
	return o.BufferSize
}

type OutputTCPSocketConfig struct {
	Host       string
	Port       uint16
	BufferSize uint
}

func (o *OutputTCPSocketConfig) Validate() *ValidationError {
	err := NewValidationError()
	if o.Host == "" {
		err.Add(errors.New("Host must not be empty"))
	}
	return err.Err()
}

func (o *OutputTCPSocketConfig) GetAddress() string {
	host := o.Host
	port := o.Port
	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = 10053
	}
	return host + ":" + strconv.Itoa(int(port))
}

func (o *OutputTCPSocketConfig) GetBufferSize() uint {
	if o.BufferSize == 0 {
		return OutputBufferSize
	}
	return o.BufferSize
}

type OutputFluentConfig struct {
	Host       string
	Tag        string
	Port       uint16
	IPv4Mask   uint8
	IPv6Mask   uint8
	BufferSize uint
}

func (o *OutputFluentConfig) Validate() *ValidationError {
	err := NewValidationError()
	if o.Host == "" {
		err.Add(errors.New("Host must not be empty"))
	}
	if o.Tag == "" {
		err.Add(errors.New("Tag must not be empty"))
	} else {
		r := regexp.MustCompile(`^[a-z0-9_]+$`)
		labels := strings.Split(o.Tag, ".")
		for _, label := range labels {
			if r.MatchString(label) {
				err.Add(errors.New("Tag characters must only include lower-case alphabets, digits underscore and dot"))
				break
			}
		}
		if o.Tag[0] == '.' {
			err.Add(errors.New("First part of a tag is empty"))
		}
		if o.Tag[len(o.Tag)-1] == '.' {
			err.Add(errors.New("Last part of a tag is empty"))
		}
		if o.IPv4Mask != 0 {
			if o.IPv4Mask > 32 {
				err.Add(errors.New("IPv4Mask must include range 0 to 32"))
			}
		}
		if o.IPv6Mask != 0 {
			if o.IPv6Mask > 128 {
				err.Add(errors.New("IPv4Mask must include range 0 to 128"))
			}
		}
	}
	return err.Err()
}

func (o *OutputFluentConfig) GetHost() string {
	return o.Host
}

func (o *OutputFluentConfig) GetTag() string {
	return o.Tag
}

func (o *OutputFluentConfig) GetPort() int {
	if o.Port == 0 {
		return 24224
	}
	return int(o.Port)
}

func (o *OutputFluentConfig) GetIPv4Mask() int {
	if o.IPv4Mask == 0 {
		return 24
	}
	return int(o.IPv4Mask)
}

func (o *OutputFluentConfig) GetIPv6Mask() int {
	if o.IPv4Mask == 0 {
		return 48
	}
	return int(o.IPv6Mask)
}

func (o *OutputFluentConfig) GetBufferSize() uint {
	if o.BufferSize == 0 {
		return OutputBufferSize
	}
	return o.BufferSize
}
