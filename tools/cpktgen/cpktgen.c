
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
		.sin_port = htons(60000),
		.sin_addr = {
			.s_addr = addr
		}
	};

	char * buffer = malloc(1280);
	if (buffer == NULL) {
		perror("Unable to allocate memory");
		exit(1);
	}

	while (1) {
		ssize_t s = sendto(dst, buffer, 1280, 0, (struct sockaddr*)&sockaddr,
				sizeof(sockaddr));
		if (s < 0) {
			perror("Unable to send datagram");
			exit(1);
		}
	}

	return 0;
}
