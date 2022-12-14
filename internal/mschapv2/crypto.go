package mschapv2

import (
	"crypto"
	_ "crypto/sha1"
)

var generateAuthenticatorResponseMagic1 = []byte{0x4D, 0x61, 0x67, 0x69, 0x63, 0x20, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x20, 0x74, 0x6F, 0x20, 0x63, 0x6C, 0x69, 0x65,
	0x6E, 0x74, 0x20, 0x73, 0x69, 0x67, 0x6E, 0x69, 0x6E, 0x67,
	0x20, 0x63, 0x6F, 0x6E, 0x73, 0x74, 0x61, 0x6E, 0x74}

var generateAuthenticatorResponseMagic2 = []byte{0x50, 0x61, 0x64, 0x20, 0x74, 0x6F, 0x20, 0x6D, 0x61, 0x6B,
	0x65, 0x20, 0x69, 0x74, 0x20, 0x64, 0x6F, 0x20, 0x6D, 0x6F,
	0x72, 0x65, 0x20, 0x74, 0x68, 0x61, 0x6E, 0x20, 0x6F, 0x6E,
	0x65, 0x20, 0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6F,
	0x6E}

// RFC2759 section 8.7 page 9
// return binary form of the result
// in the hash example password is utf-16 encode.
//    Password = "clientPass" 密码的输入格式比较奇怪.utf16?
// 			   = 63 00 6C 00 69 00 65 00 6E 00
// 				 74 00 50 00 61 00 73 00 73 00
func GenerateAuthenticatorResponse(Password []byte, NTResponse [24]byte, PeerChallenge [16]byte,
	AuthenticatorChallenge [16]byte, UserName []byte) [20]byte {
	//TODO 解决中文密码的问题
	utf16Password := make([]byte, len(Password)*2)
	for i := range Password {
		utf16Password[i*2] = Password[i]
	}
	/*
	 * Hash the password with MD4
	 */

	PasswordHash := md4(utf16Password) //ntPasswordHash

	/*
	 * Now hash the hash
	 */

	PasswordHashHash := md4(PasswordHash) //hashNtPasswordHash

	h := crypto.SHA1.New()
	h.Write(PasswordHashHash)
	h.Write(NTResponse[:])
	h.Write(generateAuthenticatorResponseMagic1)
	Digest := h.Sum(nil)

	Challenge := challengeHash(PeerChallenge, AuthenticatorChallenge, UserName)

	h = crypto.SHA1.New()
	h.Write(Digest)
	h.Write(Challenge)
	h.Write(generateAuthenticatorResponseMagic2)
	Digest = h.Sum(nil)

	out := [20]byte{}
	copy(out[:], Digest)
	return out
}

// it is md4 hash
func md4(Password []byte) []byte {
	h := crypto.MD4.New()
	h.Write(Password)
	return h.Sum(nil)
}

func challengeHash(PeerChallenge [16]byte, AuthenticatorChallenge [16]byte, UserName []byte) []byte {
	h := crypto.SHA1.New()
	h.Write(PeerChallenge[:])
	h.Write(AuthenticatorChallenge[:])
	h.Write(UserName)
	Digest := h.Sum(nil)

	Challenge := make([]byte, 8)
	copy(Challenge, Digest[:8])
	return Challenge
}

var msCHAPV2GetSendAndRecvKeySHSpad1 = [40]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

var msCHAPV2GetSendAndRecvKeySHSpad2 = [40]byte{0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2,
	0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2,
	0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2,
	0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2}

var msCHAPV2GetSendAndRecvKeyMagic1 = [27]byte{0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x74,
	0x68, 0x65, 0x20, 0x4d, 0x50, 0x50, 0x45, 0x20, 0x4d,
	0x61, 0x73, 0x74, 0x65, 0x72, 0x20, 0x4b, 0x65, 0x79}

var msCHAPV2GetSendAndRecvKeyMagic2 = [84]byte{0x4f, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x20, 0x73, 0x69, 0x64, 0x65, 0x2c, 0x20,
	0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68,
	0x65, 0x20, 0x73, 0x65, 0x6e, 0x64, 0x20, 0x6b, 0x65, 0x79,
	0x3b, 0x20, 0x6f, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x20, 0x73, 0x69, 0x64, 0x65,
	0x2c, 0x20, 0x69, 0x74, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68,
	0x65, 0x20, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x20,
	0x6b, 0x65, 0x79, 0x2e}

var msCHAPV2GetSendAndRecvKeyMagic3 = [84]byte{0x4f, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x20, 0x73, 0x69, 0x64, 0x65, 0x2c, 0x20,
	0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68,
	0x65, 0x20, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x20,
	0x6b, 0x65, 0x79, 0x3b, 0x20, 0x6f, 0x6e, 0x20, 0x74, 0x68,
	0x65, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x20, 0x73,
	0x69, 0x64, 0x65, 0x2c, 0x20, 0x69, 0x74, 0x20, 0x69, 0x73,
	0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x6e, 0x64, 0x20,
	0x6b, 0x65, 0x79, 0x2e}

func msCHAPV2GetSendAndRecvKeyGetMasterKey(passwordHashHash []byte, NTResponse [24]byte) (masterkey [16]byte) {
	// Secure Hash Standard Federal Information Processing Standards Publication 180-1 -> SHA1
	h := crypto.SHA1.New()
	h.Write(passwordHashHash)
	h.Write(NTResponse[:])
	h.Write(msCHAPV2GetSendAndRecvKeyMagic1[:])
	Digest := h.Sum(nil)
	copy(masterkey[:], Digest[:16])
	return
}

func msCHAPV2GetSendAndRecvKeyGetAsymetricStartKey(masterkey [16]byte, sessionKeyLength int, IsSend bool, IsServer bool) (sessionKey []byte) {
	var s [84]byte
	if IsSend {
		if IsServer {
			s = msCHAPV2GetSendAndRecvKeyMagic3
		} else {
			s = msCHAPV2GetSendAndRecvKeyMagic2
		}
	} else {
		if IsServer {
			s = msCHAPV2GetSendAndRecvKeyMagic2
		} else {
			s = msCHAPV2GetSendAndRecvKeyMagic3
		}
	}
	h := crypto.SHA1.New()
	h.Write(masterkey[:])
	h.Write(msCHAPV2GetSendAndRecvKeySHSpad1[:])
	h.Write(s[:])
	h.Write(msCHAPV2GetSendAndRecvKeySHSpad2[:])
	Digest := h.Sum(nil)
	sessionKey = Digest[:sessionKeyLength]
	return sessionKey
}

// http://www.ietf.org/rfc/rfc3079.txt
//    Password = "clientPass"
// 			   = 63 00 6C 00 69 00 65 00 6E 00
// 				 74 00 50 00 61 00 73 00 73 00
func MsCHAPV2GetSendAndRecvKey(Password []byte, NTResponse [24]byte) (sendKey []byte, recvKey []byte) {
	utf16Password := make([]byte, len(Password)*2)
	for i := range Password {
		utf16Password[i*2] = Password[i]
	}
	PasswordHash := md4(utf16Password)
	PasswordHashHash := md4(PasswordHash)

	masterKey := msCHAPV2GetSendAndRecvKeyGetMasterKey(PasswordHashHash, NTResponse)
	MasterSendKey := msCHAPV2GetSendAndRecvKeyGetAsymetricStartKey(masterKey, 16, true, true)
	MasterReceiveKey := msCHAPV2GetSendAndRecvKeyGetAsymetricStartKey(masterKey, 16, false, true)

	return MasterSendKey, MasterReceiveKey
}
