package html

import (
	"strings"

	"github.com/9d4/wadoh/devices"
)

// devicePhone get phone from a Device
func devicePhone (d *devices.Device) string {
    s := strings.SplitN(d.ID, ":", 2)
    if len(s) == 2 {
        return s[0]
    }
    return d.ID
}
