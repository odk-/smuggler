FROM scratch
ADD smuggler /
EXPOSE 8080
CMD ["./smuggler"]