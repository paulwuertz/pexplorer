### Regression test firmware

To ensure things are not broken and get more acurate over time some test firmware is included.

| Firmware                                                        | Test build file                          | CPU ARCH          | CPU Model       | OS       | Compiler |
| --------------------------------------------------------------- | ---------------------------------------- | ----------------- | --------------- | -------- | -------- |
| [ZSWatch](https://github.com/ZSWatch/ZSWatch)                   | zswatch_nrf5340_07.elf                   | Arm Cortex-M33    | nRF5340         | Zephyr   | GCC      |
| [IronOS](https://github.com/Ralim/IronOS)                       | Pinecilv2_EN_v2_23.elf                   | RISC-V SiFive E24 | Bouffalo BL-706 | FreeRTOS | GCC      |
| [IronOS](https://github.com/Ralim/IronOS)                       | Pinecilv1_EN_v2_23.elf                   | RISC-V RV32IMAC   | GD32VF103       | FreeRTOS | GCC      |
| [CANnectivity](github.com/CANnectivity/cannectivity)            | zephyr_cannectivity_12_llvm_lpc55s16.elf | Arm Cortex-M33    | NXP lpc55s16    | Zephyr   | LLVM     |
| [CANnectivity](github.com/CANnectivity/cannectivity)            | zephyr_cannectivity_13_llvm_lpc55s16.elf | Arm Cortex-M33    | NXP lpc55s16    | Zephyr   | LLVM     |
| [CANnectivity](github.com/CANnectivity/cannectivity)            | zephyr_cannectivity_13_gcc_lpc55s16.elf  | Arm Cortex-M33    | NXP lpc55s16    | Zephyr   | GCC      |
| [CANnectivity](github.com/CANnectivity/cannectivity)            | zephyr_cannectivity_12_gcc_lpc55s16.elf  | Arm Cortex-M33    | NXP lpc55s16    | Zephyr   | GCC      |
| [Prusa Buddy](https://github.com/prusa3d/Prusa-Firmware-Buddy/) | prusa_buddy_boot_64.elf                  | Arm Cortex-M4     | STM32F429VI     | FreeRTOS | GCC      |

#### Tests

`go test -bench=. ./cmd/fw_jsonreport/`
