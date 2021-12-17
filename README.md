# polar-serial-to-keyboard
Serial to Keyboard converter for Polar's card readers

## Configuration

This program is intended to be run as a background process. 
As such, it doesn't have any kind of user-face interface and has to use a configuration file to choose which port to connect to on runtime.

Configuration is done through the `config.json` file with two possible fields : 
- `port` : specifies which port to connect to;
- `deviceName` : specifies the device name to which to connect to.
Note that only one of the two parameter is taken into account, and `deviceName` has priority over `port`.
A utility program (`show-serial-devices`) is shipped along to help identify the device to connect to.

## Usage

Once the program is properly configured, it should appear in the taskbar tray. If it isn't the case, please verify check the `logs.log` file and check your config.

## Why is this even needed

YSoft's card readers' integrated USB keyboard mode emulates key strokes on EN/US keyboard layout, which doesn't work with Polar's terminals, which are set to FR.
