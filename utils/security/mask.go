package security

// MaskPhoneNumber 掩码国内 11 位手机号
func MaskPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) < 11 {
		return phoneNumber
	}
	return phoneNumber[:3] + "****" + phoneNumber[7:]
}
