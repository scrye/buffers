
#define _GNU_SOURCE

#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <netinet/udp.h>

int main(int argc, char * argv[]) {
	int rc;

	int dst = socket(AF_INET, SOCK_DGRAM, 0);
	if (dst < 0) {
		perror("Unable to create UDP socket");
		exit(1);
	}

	unsigned long addr;
	rc = inet_pton(AF_INET, "169.254.2.254", &addr);
	if (rc <= 0) {
		perror("Unable to parse address");
		exit(1);
	}

	struct sockaddr_in sockaddr = {
		.sin_family = AF_INET,
		.sin_port = htons(0x0800),
		.sin_addr = {
			.s_addr = addr
		}
	};

	int n = 100;
	struct mmsghdr *msgvec = calloc(n, sizeof(struct mmsghdr));
	for (int i = 0; i < n; i++) {
		msgvec[i].msg_hdr.msg_name = &sockaddr;
		msgvec[i].msg_hdr.msg_namelen = sizeof(sockaddr);

		struct iovec * iovec = malloc(sizeof(iovec));
		if (iovec == NULL) {
			perror("Unable to allocate memory");
			exit(1);
		}
		iovec->iov_base = malloc(1280);
		if (iovec->iov_base == NULL) {
			perror("Unable to allocate memory");
			exit(1);
		}
		iovec->iov_len = 1280;
		msgvec[i].msg_hdr.msg_iov = iovec;
		msgvec[i].msg_hdr.msg_iovlen = 1;
	}

	while (1) {
		rc = sendmmsg(dst, msgvec, n, 0);
		if (rc < 0) {
			perror("Unable to send datagrams");
			exit(1);
		}
	}

	return 0;
}
