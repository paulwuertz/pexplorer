package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/knightsc/gapstone"
)

func main() {

	engine, err := gapstone.New(gapstone.CS_ARCH_ARM, gapstone.CS_MODE_THUMB)
	str := "ELX/9/f/EkgSSRNK//f4/hJLG2gT8AEPDtFA9hgAAPAz+QRGR/D6/w1LHGABIgtLGmAMSLLwRPsJSxhoA2icaAlJASJP9BZzoEep8Jf9/ufcTgsI8E4LCPROCwhoBwAgZAcAIF0GAAgETwsI"
	assembly, err := base64.StdEncoding.DecodeString(str)
	fmt.Printf("%X", assembly)
	if err == nil {

		defer engine.Close()

		maj, min := engine.Version()
		log.Printf("Hello Capstone! Version: %v.%v\n", maj, min)

		// var x86Code32 = "\x8d\x4c\x32\x08\x01\xd8\x81\xc6\x34" +
		//     "\x12\x00\x00\x05\x23\x01\x00\x00\x36\x8b\x84\x91" +
		//     "\x23\x01\x00\x00\x41\x8d\x84\x39\x89\x67\x00\x00" +
		//     "\x8d\x87\x89\x67\x00\x00\xb4\xc6"

		insns, err := engine.Disasm(
			assembly, // code buffer
			0x10000,  // starting address
			0,        // insns to disassemble, 0 for all
		)

		if err == nil {
			log.Printf("Disasm:\n")
			for _, insn := range insns {
				log.Printf("0x%x:\t%s\t\t%s\n", insn.Address, insn.Mnemonic, insn.OpStr)
			}
			return
		}
		log.Fatalf("Disassembly error: %v", err)
	}
	log.Fatalf("Failed to initialize engine: %v", err)
}
