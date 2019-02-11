#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

#include <unistd.h>
#include <sys/stat.h>
#include <error.h>
#include <assert.h>
#include <sys/socket.h>
#include <netinet/in.h>

#define PORT 5522

char ** loadfile(char *filename, int *len, int sock);
int loadsock();

int main(int argc, char *argv[]) {
	if (argc == 1) {
		fprintf(stderr, "Cade o arquivo!!!\n");
		exit(1);
	}

	int length = 0;
	int sock = 0;

	if ((sock = loadsock()) <= 0) {
		exit(1);
	}

	char **words = loadfile(argv[1], &length, sock);
}

int loadsock() {
	struct sockaddr_in address;
	int sock = 0, valread;
	struct sockaddr_in serv_addr;
	
	if ((sock = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
		fprintf(stderr, "Socket create error\n");
		return -1;
	}

	memset(&serv_addr, '0', sizeof(serv_addr));

	serv_addr.sin_family = AF_INET;
	serv_addr.sin_port = htons(PORT);

	if (inet_pton(AF_INET, "172.16.95.171", &serv_addr.sin_addr) <= 0) {
		printf("Invalid Address!\n");
		return -1;
	}

	if (connect(sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0) {
		printf("Failed connection\n");
		return -1;
	}

	return sock;
}

char **loadfile(char *filename, int *len, int sock) {
	/*FILE *f = fopen(filename, "r");
	if (!f) {
		fprintf(stderr, "Can't open %s for reading\n", filename);
		return NULL;
	}*/

	int f = 0;
	if ((f = open(filename, O_RDONLY)) == -1) {
		return NULL;
	}


	char **lines = (char **)malloc(100 *sizeof(char *));
	size_t result = 0;
	int size = 1024;
	char *buf;
	int readBytes = 0;
	int oldBytes = 0;

	//while(!feof(f)) {

	while(1) {
		buf = (char*)malloc(1024 *sizeof(char *));
		//result = fread(buf, size, 1, f);
		//buf[strlen(buf)] = '\0';
		printf("%d\n", readBytes);
		//lseek(f, readBytes, SEEK_SET);
		result = read(f, buf, size);
		if (!result) {
			break;
		}
		//send(sock, buf, result, 0);
		readBytes += result;
		free(buf);
		//posix_fadvise(f, (readBytes + 1) - result, readBytes, POSIX_FADV_DONTNEED);
	}

	posix_fadvise(f, 0, 0, POSIX_FADV_DONTNEED);
	/*char **lines = (char **)malloc(100 *sizeof(char *));
	//char buf[35 * 1024];

	char *buf = (char*)malloc(1024 *sizeof(char *));
	while (fgets(buf, 1024, f)) {
	buf[strlen(buf)] = '\0';

	int slen = strlen(buf);
	printf("buf %d\n", slen);
	memset(buf, 0, 1024);
	free(buf);
	buf = (char*)malloc(1024 *sizeof(char *));
	//char * str = (char *)malloc((slen + 1) * sizeof(char));
	//strcpy(str, buf);
	}*/

	close(f);

	return lines;
}
