package rtos

import (
	"fmt"
	"slices"
	"strings"

	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func ScanForRtosFeatures(s *symbolextraction.SElfReport) {
	// TODO only zephyr for now... how to definitly detect it though...
	for _, v := range s.Variables {
		// what to get from them...
		if strings.Contains(v.Name, "dts_ord") {
			fmt.Println("found device:", v.Name, v)
		}
		// static threads created by macro
		if v.VariableType == "k_thread" {
			fmt.Println("found static thread:", v)
		}
	}
	// dynamic threads created during runtime
	var thread_create_fn symbolextraction.FunctionSymbol
	if slices.ContainsFunc(s.Functions, func(e symbolextraction.FunctionSymbol) bool {
		if e.Name == "z_impl_k_thread_create" {
			thread_create_fn = e
			return true
		}
		return false
	}) {
		thread_creating_fns := thread_create_fn.Callers
		fmt.Println("found dynamic thread:", len(thread_creating_fns))
		for i, t_addr := range thread_creating_fns {

			// struct gs_usb_data {
			// 	struct usb_dev_data common;
			// 	struct gs_usb_channel_data channels[CONFIG_USB_DEVICE_GS_USB_MAX_CHANNELS];
			// 	size_t nchannels;
			// 	struct gs_usb_ops ops;
			// 	void *user_data;
			// 	struct net_buf_pool *pool;
			// 	void *if0_str_desc;
			// 	struct k_fifo rx_fifo;
			// 	struct k_thread rx_thread;

			// #ifdef CONFIG_USB_DEVICE_GS_USB_COMPATIBILITY_MODE
			// 	uint8_t tx_buffer1[GS_USB_HOST_FRAME_MAX_SIZE];
			// #endif /* CONFIG_USB_DEVICE_GS_USB_COMPATIBILITY_MODE */
			// 	uint8_t tx_buffer2[GS_USB_HOST_FRAME_MAX_SIZE];
			// 	struct k_fifo tx_fifo;
			// 	struct k_thread tx_thread;

			// struct device {
			//         const char *name;
			//         const void *config;
			//         const void *api;
			//         struct device_state *state;
			//         void *data;
			//         struct device_ops ops;
			//         device_flags_t flags;

			// static int gs_usb_init(const struct device *dev)
			// {
			// 	struct gs_usb_data *data = dev->data;
			// 	data->common.dev = dev;
			// 	k_fifo_init(&data->rx_fifo);
			// 	k_fifo_init(&data->tx_fifo);
			// 	k_thread_create(&data->rx_thread, data->rx_stack, K_KERNEL_STACK_SIZEOF(data->rx_stack),
			// 			gs_usb_rx_thread, (void *)dev, NULL, NULL,
			// 			CONFIG_USB_DEVICE_GS_USB_RX_THREAD_PRIO, 0, K_NO_WAIT);
			// 	k_thread_name_set(&data->rx_thread, "gs_usb_rx");
			// 	k_thread_create(&data->tx_thread, data->tx_stack, K_KERNEL_STACK_SIZEOF(data->tx_stack),
			// 			gs_usb_tx_thread, (void *)dev, NULL, NULL,
			// 			CONFIG_USB_DEVICE_GS_USB_TX_THREAD_PRIO, 0, K_NO_WAIT);
			// 	k_thread_name_set(&data->tx_thread, "gs_usb_tx");
			// 	return 0;
			// }

			// k_tid_t z_impl_k_thread_create(struct k_thread *new_thread,
			//       k_thread_stack_t *stack,
			//       size_t stack_size, k_thread_entry_t entry,
			//       void *p1, void *p2, void *p3,
			//       int prio, uint32_t options, k_timeout_t delay)
			t, known := s.Addr2FnMap[*t_addr.CallFrom]
			if known {
				fmt.Println("\t#", i, "initialized in", t.Name)
			} else {
				fmt.Println("\t#", i, "initialized by unknown fn from addr", t_addr)
			}
		}
	}
}
