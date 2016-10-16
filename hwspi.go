package hwspi

import (
	"fmt"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

//HWspi is a driver for HWspi led strips
type HWspi struct {
	ClkOut    embd.DigitalPin
	ClkPin    string
	DataOut   embd.DigitalPin
	DataPin   string
	ClkFactor time.Duration
}

//Init Creates a New HWspi driver based on embd package
func (spi *HWspi) Init(ClkPin, DataPin string, ClkFactor time.Duration) (*HWspi, error) {
	spi.ClkPin = ClkPin
	spi.DataPin = DataPin
	if ClkFactor < 1 || ClkFactor > 1000 {
		spi.ClkFactor = 1
	} else {
		spi.ClkFactor = ClkFactor
	}
	embd.InitGPIO()

	clkout, err := embd.NewDigitalPin(ClkPin)
	if err != nil {
		panic(err)
	}
	spi.ClkOut = clkout
	dataout, err := embd.NewDigitalPin(DataPin)
	if err != nil {
		panic(err)
	}
	spi.DataOut = dataout

	spi.ClkOut.SetDirection(embd.Out)
	spi.DataOut.SetDirection(embd.Out)

	spi.ClkOut.Write(embd.Low)
	spi.DataOut.Write(embd.Low)

	return spi, nil
}

//GpioWriteBuffer - writes a buffer to spi
func (spi *HWspi) GpioWriteBuffer(bytes []byte) {
	for _, b := range bytes {
		spi.GpioWrite(b)
	}
}

//GpioWrite - writes a byte to spi
func (spi *HWspi) GpioWrite(b byte) {
	var i uint
	for i = 0; i < 8; i++ {
		if b&(1<<i) != 0 {
			spi.GpioWriteBit(true)
		} else {
			spi.GpioWriteBit(false)
		}
	}
}

//GpioWriteBit sends a single bit over SPI
func (spi *HWspi) GpioWriteBit(b bool) {
	if b == true {
		spi.DataOut.Write(embd.High)
		spi.gpioSynchronize()
	} else {
		spi.DataOut.Write(embd.Low)
		spi.gpioSynchronize()
	}
}

//Synchronize used to fake SPI
func (spi *HWspi) gpioSynchronize() {
	err := spi.ClkOut.Write(embd.High)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Microsecond * spi.ClkFactor)
	err = spi.ClkOut.Write(embd.Low)
	if err != nil {
		fmt.Println(err)
	}
}
