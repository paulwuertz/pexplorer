### Regression test firmware

To ensure things are not broken and get more acurate over time some test firmware is included.

| Test firmware file                       | CPU ARCH          | CPU Model       | OS       | Compiler |
|------------------------------------------|-------------------|-----------------|----------|----------|
| zswatch_nrf5340_07.elf                   | Arm Cortex-M33    | nRF5340         | Zephyr   | GCC      |
| Pinecilv2_EN_v2_23.elf                   | RISC-V SiFive E24 | Bouffalo BL-706 | FreeRTOS | GCC      |
| Pinecilv1_EN_v2_23.elf                   | RISC-V RV32IMAC   | GD32VF103       | FreeRTOS | GCC      |
| zephyr_cannectivity_12_llvm_lpc55s16.elf | Arm Cortex-M33    | NXP lpc55s16    | Zephyr   | LLVM     |
| zephyr_cannectivity_13_llvm_lpc55s16.elf | Arm Cortex-M33    | NXP lpc55s16    | Zephyr   | LLVM     |
| zephyr_cannectivity_13_gcc_lpc55s16.elf  | Arm Cortex-M33    | NXP lpc55s16    | Zephyr   | GCC      |
| zephyr_cannectivity_12_gcc_lpc55s16.elf  | Arm Cortex-M33    | NXP lpc55s16    | Zephyr   | GCC      |
| prusa_buddy_boot_64.elf                  | Arm Cortex-M4     | STM32F429VI     | FreeRTOS | GCC      |

