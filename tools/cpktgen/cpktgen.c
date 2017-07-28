
#include <stdio.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <netinet/udp.h>

int main(int argc, char * argv[]) {
	int rc;

	int dst = socket(AF_INET, SOCK_DGRAM, 0);
	if (dst < 0) {
		fprintf(stderr, "Unable to create UDP socket\n");
		exit(1);
	}

	unsigned long addr;
	rc = inet_pton(AF_INET, "169.254.2.254", &addr);
	if (rc <= 0) {
		fprintf(stderr, "Unable to parse address\n");
		exit(1);
	}

	struct sockaddr_in sockaddr = {
		.sin_family = AF_INET,
		.sin_port = 60000,
		.sin_addr = {
			.s_addr = addr
		}
	};
	rc = connect(dst, (struct sockaddr*)&sockaddr, sizeof(sockaddr));
	if (rc < 0) {
		fprintf(stderr, "Unable to connect to destination\n");
		exit(1);
	}

	char * buffer = malloc(1280);
	if (buffer == NULL) {
		fprintf(stderr, "Unable to allocate memory\n");
		exit(1);
	}

	while (1) {
		ssize_t s = send(dst, buffer, 1280, 0);	
		if (s < 0) {
			fprintf(stderr, "Unable to send datagram\n");
			exit(1);
		}
	}

	return 0;
}
