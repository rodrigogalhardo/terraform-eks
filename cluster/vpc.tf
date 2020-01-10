#
# VPC Resources
#  * VPC
#  * Subnets
#  * Internet Gateway
#  * Route Table
#

resource "aws_vpc" "cardpay" {
  cidr_block = "10.0.0.0/16"

  tags = "${
    map(
     "Name", "eks-cardpay-node-vpc",
     "kubernetes.io/cluster/${var.cluster-name}", "shared",
    )
  }"
}

resource "aws_subnet" "cardpay" {
  count = 2

  availability_zone = "${data.aws_availability_zones.available.names[count.index]}"
  cidr_block        = "10.0.${count.index}.0/24"
  vpc_id            = "${aws_vpc.cardpay.id}"

  tags = "${
    map(
     "Name", "eks-cardpay-node-subnet",
     "kubernetes.io/cluster/${var.cluster-name}", "shared",
    )
  }"
}

resource "aws_internet_gateway" "cardpay" {
  vpc_id = "${aws_vpc.cardpay.id}"

  tags {
    Name = "eks-cardpay-igw"
  }
}

resource "aws_route_table" "cardpay" {
  vpc_id = "${aws_vpc.cardpay.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.cardpay.id}"
  }
}

resource "aws_route_table_association" "cardpay" {
  count = 2

  subnet_id      = "${aws_subnet.cardpay.*.id[count.index]}"
  route_table_id = "${aws_route_table.cardpay.id}"
}
