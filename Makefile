NAME_PS = ps
NAME_CK = ck

SRC_PS = ./push-swap/main.go
SRC_CK = ./checker/main.go

all: $(NAME_PS) $(NAME_CK)

$(NAME_PS):
	go build -o $(NAME_PS) $(SRC_PS)

$(NAME_CK):
	go build -o $(NAME_CK) $(SRC_CK)

clean:
	rm -f $(NAME_PS) $(NAME_CK)

re: clean all

test:
	go test -v ./stack/...

.PHONY: all clean fclean re test
