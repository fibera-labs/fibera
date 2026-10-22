package security

import "strings"

// MaskString 保留首尾指定字符数，其余字符以 * 替代
func MaskString(value string, keepPrefix, keepSuffix int) string {
	runes := []rune(value)
	if keepPrefix < 0 {
		keepPrefix = 0
	}
	if keepSuffix < 0 {
		keepSuffix = 0
	}
	if keepPrefix+keepSuffix >= len(runes) {
		return value
	}
	return string(runes[:keepPrefix]) + strings.Repeat("*", len(runes)-keepPrefix-keepSuffix) + string(runes[len(runes)-keepSuffix:])
}

// MaskPhoneNumber 掩码国内 11 位手机号，保留前 3 位与后 4 位
func MaskPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) < 11 {
		return phoneNumber
	}
	return MaskString(phoneNumber, 3, 4)
}

// MaskEmail 掩码邮箱地址，仅保留用户名首字符与完整域名
func MaskEmail(email string) string {
	at := strings.LastIndex(email, "@")
	if at <= 0 {
		return email
	}
	localPart, domain := email[:at], email[at:]
	localRunes := []rune(localPart)
	if len(localRunes) == 1 {
		return "*" + domain
	}
	return string(localRunes[:1]) + "***" + domain
}

// MaskName 掩码姓名，仅保留首字符
func MaskName(name string) string {
	runes := []rune(name)
	if len(runes) <= 1 {
		return name
	}
	return string(runes[:1]) + strings.Repeat("*", len(runes)-1)
}

// MaskIDCard 掩码身份证号，保留前 6 位与后 4 位
func MaskIDCard(idCard string) string {
	if len(idCard) <= 10 {
		return idCard
	}
	return MaskString(idCard, 6, 4)
}

// MaskBankCard 掩码银行卡号，仅保留后 4 位
func MaskBankCard(cardNumber string) string {
	if len(cardNumber) <= 4 {
		return cardNumber
	}
	return strings.Repeat("*", len(cardNumber)-4) + cardNumber[len(cardNumber)-4:]
}
