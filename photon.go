package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

type PhotonFile struct {
	FilePath string
	Config   *PhotonFileConfig
	Layers   int
}

/*
func CreateEmptyPhotonFile(outputFilepath string) (*PhotonFile, error) {
	photonFile := &PhotonFile{
		FilePath: outputFilepath,
	}

	return photonFile, nil
}
*/

func ReadPhotonFile(filepath string) (*PhotonFile, error) {
	photonFile := &PhotonFile{
		FilePath: filepath,
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read .photon file: %v", err)
	}

	//Read config
	photonFile.Config, err = configFromBytes(data)

	return photonFile, err
}

type PhotonFileConfig struct {
	Header                         int32
	Version                        int32
	BedX                           float32
	BedY                           float32
	BedZ                           float32
	layer_height                   float32
	exposure_time                  float32
	exposure_time_bottom           float32
	off_time                       float32
	bottom_layers                  int32
	resolution_x                   int32
	resolution_y                   int32
	preview_highres_header_address int32
	layer_def_address              int32
	n_layers                       int32
	preview_lowres_header_address  int32
	print_time                     int32 //if version > 1 or padding
	projection_type                int32
	layer_levels                   int32 //if version > 1 same as anti_aliasing_level
	print_properties_address       int32 //if version > 1
	print_properties_length        int32 //if version > 1
	anti_aliasing_level            int32 //if version > 1
	light_pwm                      int16 //if version > 1
	light_pwm_bottom               int16 //if version > 1

	//Seek preview_highres_header_address
	preview_highres_resolution_x int32
	preview_highres_resolution_y int32
	preview_highres_data_address int32
	preview_highres_data_length  int32

	//Seek preview_highres_data_address
	preview_highres_data []byte //preview_highres_data_length

	//Seek preview_lowres_header_address
	preview_lowres_resolution_x int32
	preview_lowres_resolution_y int32
	preview_lowres_data_address int32
	preview_lowres_data_length  int32

	//Seek preview_lowres_data_address
	preview_lowres_data []byte //preview_lowres_data_length

	//Seek print_properties_address
	bottom_lift_distance   float32 //if version > 1
	bottom_lift_speed      float32 //if version > 1
	lifting_distance       float32 //if version > 1
	lifting_speed          float32 //if version > 1
	retract_speed          float32 //if version > 1
	volume_ml              float32 //if version > 1
	weight_g               float32 //if version > 1
	cost_dollars           float32 //if version > 1
	bottom_light_off_delay float32 //if version > 1
	light_off_delay        float32 //if version > 1
	bottom_layer_count     int32   //if version > 1
	p1                     float32 //if version > 1
	p2                     float32 //if version > 1
	p3                     float32 //if version > 1
	p4                     float32 //if version > 1
}

//Based on https://sourcegraph.com/github.com/Photonsters/PhotonFileValidator@0a91b655a5e6602546cbfe812306c3ed79cb03ee/-/blob/src/photon/file/parts/photon/PhotonFileHeader.java#L37:14

func configFromBytes(data []byte) (*PhotonFileConfig, error) {
	out := &PhotonFileConfig{}

	buf := bytes.NewReader(data) //TODO use direct file access io.Reader

	/*
		v := reflect.ValueOf(out) //TODO not use reflect to loop over config
			for i := 0; i < v.NumField(); i++ {
				log.Debug().Msgf("DEBUG decode: %#v", v.Type().Field(i).Name)
				err := binary.Read(buf, binary.LittleEndian, v.Field(i).Interface())
				if err != nil {
					return out, fmt.Errorf("failed to decode photon config %s: %v", v.Type().Field(i).Name, err)
				}
			}
	*/
	//*
	err := binary.Read(buf, binary.LittleEndian, &out.Header)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.Version)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.BedX)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.BedY)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.BedZ)
	if err != nil {
		return out, err
	}

	//Skip 3 padding
	_, err = buf.Seek(3*4, io.SeekCurrent)
	if err != nil {
		return out, err
	}

	err = binary.Read(buf, binary.LittleEndian, &out.layer_height)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.exposure_time)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.exposure_time_bottom)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.off_time)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.bottom_layers)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.resolution_x)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.resolution_y)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.preview_highres_header_address)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.layer_def_address)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.n_layers)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.preview_lowres_header_address)
	if err != nil {
		return out, err
	}

	if out.Version > 1 {
		err = binary.Read(buf, binary.LittleEndian, &out.BedZ)
		if err != nil {
			return out, err
		}
	} else {
		//Skip 1 padding
		_, err = buf.Seek(1*4, io.SeekCurrent)
		if err != nil {
			return out, err
		}
	}
	err = binary.Read(buf, binary.LittleEndian, &out.projection_type)
	if err != nil {
		return out, err
	}

	if out.Version > 1 {
		err = binary.Read(buf, binary.LittleEndian, &out.layer_levels)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.print_properties_address)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.print_properties_length)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.anti_aliasing_level)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.light_pwm)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.light_pwm_bottom)
		if err != nil {
			return out, err
		}
	}

	//Seek preview_highres_header_address
	_, err = buf.Seek(int64(out.preview_highres_header_address), io.SeekStart)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.preview_highres_resolution_x)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.preview_highres_resolution_y)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.preview_highres_data_address)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.preview_highres_data_length)
	if err != nil {
		return out, err
	}

	//Seek preview_highres_data_address
	_, err = buf.Seek(int64(out.preview_highres_data_address), io.SeekStart)
	if err != nil {
		return out, err
	}
	out.preview_highres_data = make([]byte, out.preview_highres_data_length)
	_, err = buf.Read(out.preview_highres_data) //TODO check read length
	if err != nil {
		return out, err
	}

	//Seek preview_lowres_header_address
	_, err = buf.Seek(int64(out.preview_lowres_header_address), io.SeekStart)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.preview_lowres_resolution_x)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.preview_lowres_resolution_y)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.preview_lowres_data_address)
	if err != nil {
		return out, err
	}
	err = binary.Read(buf, binary.LittleEndian, &out.preview_lowres_data_length)
	if err != nil {
		return out, err
	}

	//Seek preview_lowres_data_address
	_, err = buf.Seek(int64(out.preview_lowres_data_address), io.SeekStart)
	if err != nil {
		return out, err
	}
	out.preview_lowres_data = make([]byte, out.preview_lowres_data_length)
	_, err = buf.Read(out.preview_lowres_data) //TODO check read length
	if err != nil {
		return out, err
	}

	if out.Version > 1 {
		//Seek print_properties_address
		_, err = buf.Seek(int64(out.print_properties_address), io.SeekStart)
		if err != nil {
			return out, err
		}

		err = binary.Read(buf, binary.LittleEndian, &out.bottom_lift_distance)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.bottom_lift_speed)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.lifting_distance)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.lifting_speed)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.retract_speed)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.volume_ml)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.weight_g)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.cost_dollars)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.bottom_light_off_delay)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.light_off_delay)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.bottom_layer_count)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.p1)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.p2)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.p3)
		if err != nil {
			return out, err
		}
		err = binary.Read(buf, binary.LittleEndian, &out.p4)
		if err != nil {
			return out, err
		}
	}

	return out, nil
}
