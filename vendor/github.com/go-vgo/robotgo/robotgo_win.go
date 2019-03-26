// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

// +build windows

package robotgo

import (
	"syscall"

	"github.com/lxn/win"
)

// FindWindow find window hwnd by name
func FindWindow(str string) win.HWND {
	hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr(str))

	return hwnd
}

// GetHWND get foreground window hwnd
func GetHWND() win.HWND {
	hwnd := win.GetForegroundWindow()

	return hwnd
}
