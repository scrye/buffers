#include <fcntl.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/ioctl.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <linux/if.h>
#include <linux/if_tun.h>

int tun_alloc(char *dev) {
	struct ifreq ifr;
	int fd, rc;

	fd = open("/dev/net/tun", O_RDWR);
	if (fd < 0) {
		perror("open");
		exit(1);
	}

	memset(&ifr, 0, sizeof(ifr));
	ifr.ifr_flags = IFF_TUN;
	if (*dev) {
		strncpy(ifr.ifr_name, dev, IFNAMSIZ);
	}

	rc = ioctl(fd, TUNSETIFF, (void *) &ifr);
	if (rc < 0) {
		close(fd);
		perror("ioctl");
		exit(1);
	}
	//strcpy(dev, ifr.ifr_name);
	return fd;
}

int main(int argc, char * argv[]) {
	int fd, n;
	char buf[1500];
	fd = tun_alloc("tuntest");

	while (1) {
		n = read(fd, buf, sizeof(buf));
		if (n < 0) {
			perror("read");
			exit(1);
		}
		printf("Read %d bytes from tuntest\n", n);
	}

	return 0;
}
