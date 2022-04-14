/*
Copyright 2022 Agus Imam Fauzi <agus7fauzi@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

type Options struct {
	Name    string
	Address string
	Version string
}

func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}
