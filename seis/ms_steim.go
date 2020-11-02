package seis

import (
	"encoding/binary"
	"fmt"
)

func getNibble(word []byte, index int) uint8 {
	b := word[index/4]                //Which byte we want from within the word (4 bytes per word)
	var res uint8                     //value
	i := index % 4                    //which nibble we want from within the byte (4 nibbles per byte)
	res = b & (0x3 << uint8((3-i)*2)) //0x3=00000011 create and apply the correct mask e.g. i=1 mask=00110000
	res = res >> uint8((3-i)*2)       //shift the masked value fully to the right
	return res
}

//value must be 0, 1, 2 or 3, the nibble must not have been previously set
func writeNibble(word []byte, index int, value uint8) {
	b := word[index/4]
	i := index % 4
	b = b ^ (value << uint8((3-i)*2)) //set the bits
	word[index/4] = b
}

/*
	Takes v: an integer where only the first numbits bits are used to represent the number and returns an int32
*/
func uintVarToInt32(v uint32, numbits uint8) int32 {
	neg := (v & (0x1 << (numbits - 1))) != 0 //check positive/negative
	if neg {                                 //2s complement
		v = v ^ ((1 << (numbits)) - 1) //flip all the bits
		v = v + 1                      //add 1 - positive nbit number
		v = -v                         //get the negative - this gives us a proper negative int32
	}
	return int32(v)
}

func int32ToUintVar(i int32, numbits uint8) uint32 {
	neg := i < 0
	if neg {
		i = -i //get the positive - this gives us a positive n bit int
		i = i - 1
		i = i ^ ((1 << (numbits)) - 1) //flip all the bits
	}
	return uint32(i)
}

/** -- steim encoding
//return the pack level for a record
func calculatePackLevelSteim(version int, d int32) (packLvl int, err error) {
	if version == 2 {
		if d <= 7 && d >= -8 { //7 x 4 bit differences
			packLvl = 7
		} else if d <= 15 && d >= -16 { //6 x 5 bit differences
			packLvl = 6
		} else if d <= 31 && d >= -32 { //5 x 6 bit differences
			packLvl = 5
		} else if d <= 127 && d >= -128 { //4 x 8 bit differences
			packLvl = 4
		} else if d <= 511 && d >= -512 { //3 x 10 bit differences
			packLvl = 3
		} else if d <= 16383 && d >= -16384 { //2 x 15 bit differences
			packLvl = 2
		} else if d <= 536870911 && d >= -536870912 { //1 x 30 bit differences
			packLvl = 1
		} else {
			return -1, fmt.Errorf("steim2: unable to encode difference value of %v", d)
		}
	} else if version == 1 {
		if d <= 127 && d >= -128 { //4 x 8 bit differences
			packLvl = 4
		} else if d <= 32767 && d >= -32768 { //2 x 16 bit differences
			packLvl = 2
		} else { //1 x 32 bit differences
			packLvl = 1
		}
		return
	}
	return
}

//get the correct nib, dnib and number of bits for the specified packLvl
func getSteimPackDetails(version, packLvl int) (nib, dnib, numbits uint8) {
	if version == 2 {
		switch packLvl {
		case 7:
			nib = 3
			dnib = 2
			numbits = 4
		case 6:
			nib = 3
			dnib = 1
			numbits = 5
		case 5:
			nib = 3
			dnib = 0
			numbits = 6
		case 4:
			nib = 1
			dnib = 0
			numbits = 8
		case 3:
			nib = 2
			dnib = 3
			numbits = 10
		case 2:
			nib = 2
			dnib = 2
			numbits = 15
		case 1:
			nib = 2
			dnib = 1
			numbits = 30
		}
	} else if version == 1 {
		switch packLvl {
		case 4:
			nib = 1
			numbits = 8
		case 2:
			nib = 2
			numbits = 16
		case 1:
			nib = 3
			numbits = 32
		}
		return
	}
	return
}

type steimOutput struct {
	EncodedData []byte
	FrameCount  uint8
	WordOrder   uint8 //0 = Little Endian, 1 = Big Endian
}

func encodeSteim(version int, data []int32) (steimOutput, error) {
	so := steimOutput{WordOrder: 1} //Only supports Big Endian

	//calculate a worst-case size for the output so we can size the array (we will trim afterwards)
	reqLength := (len(data) + 1 + 2) * 4        //4bytes for each record + 1 for the 'firstDiff' + 2 for the integration constants
	reqLength += ((len(data) + 1 + 2) / 15) * 4 //4bytes for w0 of each frame + 1 for a partially filled frame (each frame only has 15 diff words; w0 is the header)
	buf := make([]byte, reqLength+128)

	var bPointer int //Our place in the buffer (this value / 4 is the number of words we've written)
	var w0 []byte    //Where we write nibbles

	//precalculate the diffs - this probably isn't necessary but it makes the logic a bit simpler
	diffs := make([]int32, 0, len(data)+1)
	lastD := data[0] //this'll cause a 0 to be appended: the 'firstDiff' TODO: Is setting firstDiff to 0 ok?
	for _, d := range data {
		diffs = append(diffs, d-lastD)
		lastD = d
	}
	var index int //Our position in the diff

	for { //loops once for each word written
		if index == len(diffs) {
			break //we're done - no more data to write
		}
		if bPointer%64 == 0 { //create a w0 and start a new frame
			w0 = buf[bPointer : bPointer+4]
			bPointer += 4   //Keep track of where we are
			so.FrameCount++ //This is the start of a new frame
		}

		//special values
		if index == 0 {
			binary.BigEndian.PutUint32(buf[bPointer:bPointer+4], uint32(data[0])) //set the start value (forward integration constant)
			bPointer += 4

			binary.BigEndian.PutUint32(buf[bPointer:bPointer+4], uint32(data[len(data)-1])) //set the end value (reverse integration constant)
			bPointer += 4
		}

		dPack := make([]int32, 0, 7) //diffs to pack into a single word, max capacity is 7*4bit differences
		packLvl := 8                 //default the value to 8 so it can be reduced
		for {                        //loop and fill dPack till it's ready to be packed
			if index == len(diffs) { //we've reached the end of the data
				break
			}

			npackLvl, err := calculatePackLevelSteim(version, diffs[index]) //get the pack level for the next number
			if err != nil {
				return so, err
			}

			if npackLvl < packLvl { //If the new packing level is less dense than the previous one
				packLvl = npackLvl //update it
			}

			if packLvl <= len(dPack)+1 { //If the new number is too large to pack this time or is the final number to pack
				if packLvl == len(dPack)+1 { //If the new number is the final one to pack
					dPack = append(dPack, diffs[index]) //include it
					index++                             //move on to the next number
					break                               //do the pack outside the loop
				}
				break //do the pack without including the number and leave it for the next loop
			}
			dPack = append(dPack, diffs[index]) //include the new number
			index++                             //move on to the next number and loop again
		}

		nib, dnib, numbits := getSteimPackDetails(version, len(dPack)) //details for packing

		writeNibble(w0, (bPointer%64)/4, nib)

		var wint uint32 //an empty uint32 to pack the diffs into
		for j, dp := range dPack {
			uint := int32ToUintVar(dp, numbits)
			wint = wint ^ (uint << uint32(((len(dPack)-j)-1)*int(numbits))) //shift the value the correct number of bits over
		}

		w := make([]byte, 4)
		binary.BigEndian.PutUint32(w, wint)

		if version == 2 && dnib < 4 {
			writeNibble(w, 0, dnib) //write the dnib
		}

		copy(buf[bPointer:bPointer+4], w) //copy the packed value into the buffer
		bPointer += 4
	}

	so.EncodedData = buf[:int(so.FrameCount)*64] //trim the buffer to the final size TODO: Could we trim a sub frame?
	return so, nil
}
**/

func applyDifferencesFromWord(w []byte, numdiffs int, diffbits uint32, d []int32) []int32 {
	mask := (1 << diffbits) - 1
	wint := binary.BigEndian.Uint32(w)

	for i := numdiffs - 1; i >= 0; i-- {
		intn := wint & uint32(mask<<(uint32(i)*diffbits)) //apply a mask over the correct bits
		intn = intn >> (uint32(i) * diffbits)             //shift the masked value fully to the right

		diff := uintVarToInt32(intn, uint8(diffbits)) //convert diffbits bit int to int32

		d = append(d, d[len(d)-1]+diff)
	}

	return d
}

func decodeSteim(version int, raw []byte, wordOrder, frameCount uint8, expectedSamples uint16) ([]int32, error) {
	d := make([]int32, 0, expectedSamples)

	if wordOrder == 0 {
		return d, fmt.Errorf("steim%v: no support for little endian", version)
	}

	//Word 1 and 2 contain x0 and xn: the uncompressed initial and final quantities (word 0 contains nibs)
	frame0 := raw[0:64]
	start := int32(binary.BigEndian.Uint32(frame0[4:8]))
	end := int32(binary.BigEndian.Uint32(frame0[8:12]))

	d = append(d, start)

	for f := 0; f < int(frameCount); f++ {
		//Each frame is 64bytes
		frameBytes := raw[64*f : 64*(f+1)]

		//The first word (w0) of the frames contains 16 'nibbles' (2bit codes)
		w0 := frameBytes[:4]

		//Each nibble describes the encoding on a 32bit word of the frame
		//0 is w0 so we can skip that
		for w := 1; w < 16; w++ {
			nib := getNibble(w0, w) //get relevant nib

			if nib == 0 { //A 0 nib represents a non-data or header word
				continue
			}

			wb := frameBytes[w*4 : (w+1)*4] //Get word w

			//dnib is the second part of the encoding specifier and is the first nib of the word (not used in steim1)
			var dnib uint8
			if version == 2 && nib != 1 { //nib == 1 indicates no dnib (steim1 8bit encoding
				dnib = getNibble(wb, 0)
			}

			skipFirstDiff := 0
			if w == 3 && f == 0 { //libmseed seems to skip the first diff? TODO: Understand
				skipFirstDiff = 1
			}

			if version == 2 {
				switch nib {
				case 1: //4 8bit differences
					d = applyDifferencesFromWord(wb, 4-skipFirstDiff, 8, d)
				case 2:
					switch dnib {
					case 0:
						return d, fmt.Errorf("steim%v: nib 10 dnib 00 is an illegal configuration @ frame %v word %v", version, f, w)
					case 1:
						d = applyDifferencesFromWord(wb, 1-skipFirstDiff, 30, d)
					case 2:
						d = applyDifferencesFromWord(wb, 2-skipFirstDiff, 15, d)
					case 3:
						d = applyDifferencesFromWord(wb, 3-skipFirstDiff, 10, d)
					}
				case 3:
					switch dnib {
					case 0: //5 6bit differences
						d = applyDifferencesFromWord(wb, 5-skipFirstDiff, 6, d)
					case 1: //6 5bit differences
						d = applyDifferencesFromWord(wb, 6-skipFirstDiff, 5, d)
					case 2: //7 4bit differences
						d = applyDifferencesFromWord(wb, 7-skipFirstDiff, 4, d)
					case 3:
						return d, fmt.Errorf("steim%v: nib 11 dnib 11 is an illegal configuration @ frame %v word %v", version, f, w)
					}
				}
			} else if version == 1 {
				switch nib {
				case 1: //4 8bit differences
					d = applyDifferencesFromWord(wb, 4-skipFirstDiff, 8, d)
				case 2: //2 16bit differences
					d = applyDifferencesFromWord(wb, 2-skipFirstDiff, 16, d)
				case 3: //1 32bit difference
					d = applyDifferencesFromWord(wb, 1-skipFirstDiff, 32, d)
				}
			}
		}
	}

	if d[len(d)-1] != end {
		return d, fmt.Errorf("steim%v: final decompressed value did not equal reverse integration value got: %v expected %v", version, d[len(d)-1], end)
	}

	return d, nil
}
