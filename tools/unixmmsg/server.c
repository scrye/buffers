
#define _GNU_SOURCE

#include <fcntl.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/ioctl.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>

int main(int argc, char * argv[]) {
	int rc;

	int sk = socket(AF_UNIX, SOCK_DGRAM, 0);
	if (sk < 0) {
		perror("Unable to create UNIX socket");
		exit(1);
	}

	const char * address = "/tmp/unix.42";
	struct sockaddr_un sockaddr;
	sockaddr.sun_family = AF_UNIX;
	strcpy(sockaddr.sun_path, address);

	rc = bind(sk, &sockaddr, sizeof(sockaddr));
	if (rc != 0) {
		perror("Unable to bind to UNIX socket");
		exit(1);
	}

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
		rc = recvmmsg(sk, msgvec, n, 0, NULL);
		if (rc < 0) {
			perror("Unable to receive datagrams");
			exit(1);
		}
		printf("Read %d messages.\n", rc);

	}
}