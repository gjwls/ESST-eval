#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <time.h>

#define TARGET_IP "8.8.8.8" 
#define TARGET_PORT 53       
#define MESSAGE "test"       

int main() {
    struct sockaddr_in target_addr;
    char buffer[] = MESSAGE;
    int i = 1;
    long sum = 0;

    memset(&target_addr, 0, sizeof(target_addr));
    target_addr.sin_family = AF_INET;
    target_addr.sin_port = htons(TARGET_PORT);
    inet_pton(AF_INET, TARGET_IP, &target_addr.sin_addr);
    for(int j = 0; j < 10; j++){
        while (i < 100001) {
            struct timespec start, end; 
            clock_gettime(CLOCK_MONOTONIC, &start); 

            int sockfd = socket(AF_INET, SOCK_DGRAM, 0);
            if (sockfd < 0) {
                perror("failed to create socket");
                continue;
            }

            ssize_t sent_bytes = sendto(sockfd, buffer, strlen(buffer), 0,
                                        (struct sockaddr*)&target_addr, sizeof(target_addr));

            clock_gettime(CLOCK_MONOTONIC, &end); 

            if (sent_bytes < 0) {
                perror("failed to send message");
            } else {
                long elapsed_time = (end.tv_sec - start.tv_sec) * 1000000000L +
                                    (end.tv_nsec - start.tv_nsec);
                sum += elapsed_time;
            }

            close(sockfd);
            i++;
        }
        printf("%ld\n", sum);
        sum = 0;
        i = 1;
    }
    
    printf("Total transfer time: %ld\n", sum);

    return 0;
}