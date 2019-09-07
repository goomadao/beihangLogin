package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func GetEncodedPassword(d, k string) string {
	return hmac_md5(d, k)
}

func hmac_md5(d, k string) string {
	hmac := hmac.New(md5.New, []byte(k))
	hmac.Write([]byte(d))
	return hex.EncodeToString(hmac.Sum([]byte("")))
}

func GetEncodedInfo(info, token string) string {
	// fmt.Println(info)
	// fmt.Println(token)
	// fmt.Println(trashBase64(xEncode(info, token)))
	return "{SRBX1}" + trashBase64(xEncode(info, token))
}

func xEncode(info, token string) string {
	s := func(a string, b bool) []uint32 {
		c := len(a)
		v := make([]uint32, (c+3)/4)
		for i := 0; i < c; i += 4 {
			// v[i>>2] = uint32(a[i]) | uint32(a[i+1])<<8 | uint32(a[i+2])<<16 | uint32(a[i+3])<<24
			for j := 0; j < 4 && j+i < c; j++ {
				v[i>>2] |= (uint32(a[i+j]) << uint(8*j))
			}
		}
		if b {
			v = append(v, uint32(c))
		}
		return v
	}
	l := func(a []uint32, b bool) string {
		d := len(a)
		c := (d - 1) << 2
		if b {
			m := a[d-1]
			if int(m) < c-3 || int(m) > c {
				return ""
			}
			c = int(m)
		}
		var tempA bytes.Buffer
		for i := 0; i < d; i++ {
			tempA.Write([]byte{byte(a[i] & 0xff), byte(a[i] >> 8 & 0xff), byte(a[i] >> 16 & 0xff), byte(a[i] >> 24 & 0xff)})
		}
		// var str string
		// if b {
		// 	str = tempA.String()[:c]
		// } else {
		// 	str = tempA.String()
		// }
		str := tempA.String()
		if b {
			return str[0:c]
		}
		return str
	}

	if len(info) == 0 {
		return ""
	}
	v := s(info, true)
	k := s(token, false)
	n := len(v) - 1
	z := v[n]
	y := v[0]
	d := uint32(0)
	for q := 6 + 52/(n+1); q > 0; q-- {
		d += 0x9E3779B9
		e := (d >> 2) & 3
		for p := 0; p <= n; p++ {
			if p == n {
				y = v[0]
			} else {
				y = v[p+1]
			}
			m := (z >> 5) ^ (y << 2)
			m += (y >> 3) ^ (z << 4) ^ (d ^ y)
			m += k[(p&3)^int(e)] ^ z
			v[p] += m
			z = v[p]
		}
	}
	return l(v, false)
}

func trashBase64(t string) string {
	const base64N = "LVoJPiCN2R8G90yg+hmFHuacZ1OWMnrsSTXkYpUq/3dlbfKwv6xztjI7DeBE45QA"
	a := len(t)
	len := a / 3 * 4
	if a%3 != 0 {
		len += 4
	}
	u := make([]byte, len)
	r := byte('=')
	ui := 0
	for o := 0; o < a; o += 3 {
		var p [3]byte
		p[2] = t[o]
		if o+1 < a {
			p[1] = t[o+1]
		} else {
			p[1] = 0
		}
		if o+2 < a {
			p[0] = t[o+2]
		} else {
			p[0] = 0
		}
		h := int(p[2])<<16 | int(p[1])<<8 | int(p[0])
		for i := 0; i < 4; i++ {
			if o*8+i*6 > a*8 {
				u[ui] = r
			} else {
				u[ui] = base64N[h>>uint(6*(3-i))&0x3F]
			}
			ui++
		}
	}
	return string(u[:len])
}

func GetEncodedChkstr(s string) string {
	sha := sha1.New()
	sha.Write([]byte(s))
	return hex.EncodeToString(sha.Sum(nil))
}
