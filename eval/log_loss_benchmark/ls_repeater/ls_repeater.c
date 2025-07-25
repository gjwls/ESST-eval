#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char *argv[]) {
    if (argc < 2) {
        fprintf(stderr, "Usage: %s <Latency(ms)>\n", argv[0]);
        return 1;
    }

    int lat = atoi(argv[1]); 
    if (lat < 0) {
        fprintf(stderr, "Invalid latency values\n");
        return 1;
    }

    for (int i = 0; i < 100; i++) {
        printf("Executing execve(\"/usr/bin/ls\", [\"ls\"], ...)\n");
        fflush(stdout); 

        char *args[] = {"/usr/bin/ls", "ls", NULL};
        execve(args[0], args, NULL);

        perror("execve failed");
        fprintf(stderr, "execve(\"/usr/bin/ls\") failed\n");
        usleep(lat); 
    }

    return 0;
}
