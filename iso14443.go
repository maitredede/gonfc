package gonfc

// ISO14443CascadeUID Add cascade tags (0x88) in UID. ISO/IEC 14443-3 (6.4.4 UID contents and cascade levels)
func ISO14443CascadeUID(abtUID []byte, pbtCascadedUID []byte) int {
	szUID := len(abtUID)
	switch szUID {
	case 7:
		pbtCascadedUID[0] = 0x88
		// memcpy(pbtCascadedUID+1, abtUID, 7)
		for i := 0; i < 7; i++ {
			pbtCascadedUID[i+1] = abtUID[i]
		}
		// *pszCascadedUID = 8
		// break
		return 8
	case 10:
		pbtCascadedUID[0] = 0x88
		// memcpy(pbtCascadedUID+1, abtUID, 3)
		for i := 0; i < 3; i++ {
			pbtCascadedUID[i+1] = abtUID[i]
		}
		pbtCascadedUID[4] = 0x88
		// memcpy(pbtCascadedUID+5, abtUID+3, 7)
		for i := 0; i < 7; i++ {
			pbtCascadedUID[i+5] = abtUID[i]
		}
		// *pszCascadedUID = 12
		// break
		return 12
	case 4:
		fallthrough
	default:
		//memcpy(pbtCascadedUID, abtUID, szUID)
		for i := 0; i < szUID; i++ {
			pbtCascadedUID[i] = abtUID[i]
		}
		// *pszCascadedUID = szUID
		// break
		return szUID
	}
}
