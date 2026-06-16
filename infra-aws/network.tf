resource "aws_vpc" "aws_vpc" {
  cidr_block           = var.aws_vpc["cidr_block"]
  enable_dns_support   = var.aws_vpc["enable_dns_support"]
  enable_dns_hostnames = var.aws_vpc["enable_dns_hostnames"]
  tags = {
    Name = var.aws_vpc["tags"]["Name"]
  }
}

resource "aws_subnet" "aws_public_subnet1" {
  vpc_id            = aws_vpc.aws_vpc.id
  availability_zone = var.aws_public_subnet1["availability_zone"]
  cidr_block        = var.aws_public_subnet1["cidr_block"]
  tags = {
    Name = var.aws_public_subnet1["tags"]["Name"]
  }
  map_public_ip_on_launch = var.aws_public_subnet1["map_public_ip_on_launch"]
}

resource "aws_subnet" "aws_public_subnet2" {
  vpc_id            = aws_vpc.aws_vpc.id
  availability_zone = var.aws_public_subnet2["availability_zone"]
  cidr_block        = var.aws_public_subnet2["cidr_block"]
  tags = {
    Name = var.aws_public_subnet2["tags"]["Name"]
  }
  map_public_ip_on_launch = var.aws_public_subnet2["map_public_ip_on_launch"]

}

resource "aws_subnet" "aws_public_subnet3" {
  vpc_id            = aws_vpc.aws_vpc.id
  availability_zone = var.aws_public_subnet3["availability_zone"]
  cidr_block        = var.aws_public_subnet3["cidr_block"]
  tags = {
    Name = var.aws_public_subnet3["tags"]["Name"]
  }
  map_public_ip_on_launch = var.aws_public_subnet3["map_public_ip_on_launch"]
}

resource "aws_internet_gateway" "aws_internet_gateway" {
  vpc_id = aws_vpc.aws_vpc.id
  tags = {
    Name = var.aws_internet_gateway["tags"]["Name"]
  }
}

resource "aws_route_table" "aws_route_table_public_subnets" {
  vpc_id = aws_vpc.aws_vpc.id
  tags = {
    Name = var.aws_route_table_public_subnets["tags"]["Name"]
  }
  route {
    cidr_block = var.aws_route_table_public_subnets["route"]["cidr_block"]
    gateway_id = aws_internet_gateway.aws_internet_gateway.id
  }
}

resource "aws_route_table_association" "aws_route_table_association_public_subnets" {
  subnet_id      = aws_subnet.aws_public_subnet1.id
  route_table_id = aws_route_table.aws_route_table_public_subnets.id
}

resource "aws_route_table_association" "aws_route_table_association_public_subnets2" {
  subnet_id      = aws_subnet.aws_public_subnet2.id
  route_table_id = aws_route_table.aws_route_table_public_subnets.id
}

resource "aws_route_table_association" "aws_route_table_association_public_subnets3" {
  subnet_id      = aws_subnet.aws_public_subnet3.id
  route_table_id = aws_route_table.aws_route_table_public_subnets.id
}


### creating private subnets 

resource "aws_subnet" "aws_private_subnet1" {
  vpc_id            = aws_vpc.aws_vpc.id
  availability_zone = var.aws_private_subnet1["availability_zone"]
  cidr_block        = var.aws_private_subnet1["cidr_block"]
  tags = {
    Name = var.aws_private_subnet1["tags"]["Name"]
  }
}

resource "aws_eip" "aws_nat_eip_private_subnet1" {
  domain = var.aws_nat_eip_private_subnet1["domain"]
  tags = {
    Name = var.aws_nat_eip_private_subnet1["tags"]["Name"]
  }
}

resource "aws_nat_gateway" "aws_nat_private_subnet1" {
  allocation_id = aws_eip.aws_nat_eip_private_subnet1.id
  subnet_id     = aws_subnet.aws_public_subnet1.id
  tags = {
    Name = var.aws_nat_private_subnet1["tags"]["Name"]
  }
  depends_on = [aws_internet_gateway.aws_internet_gateway]
}

resource "aws_route_table" "aws_route_table_private_subnet1" {
  vpc_id = aws_vpc.aws_vpc.id
  tags = {
    Name = var.aws_route_table_private_subnet1["tags"]["Name"]
  }
  route {
    cidr_block     = var.aws_route_table_private_subnet1["route"]["cidr_block"]
    nat_gateway_id = aws_nat_gateway.aws_nat_private_subnet1.id
  }
}

resource "aws_route_table_association" "aws_route_table_association_private_subnet1" {
  subnet_id      = aws_subnet.aws_private_subnet1.id
  route_table_id = aws_route_table.aws_route_table_private_subnet1.id
}

resource "aws_subnet" "aws_private_subnet2" {
  vpc_id            = aws_vpc.aws_vpc.id
  availability_zone = var.aws_private_subnet2["availability_zone"]
  cidr_block        = var.aws_private_subnet2["cidr_block"]
  tags = {
    Name = var.aws_private_subnet2["tags"]["Name"]
  }
}

resource "aws_eip" "aws_nat_eip_private_subnet2" {
  domain = var.aws_nat_eip_private_subnet2["domain"]
  tags = {
    Name = var.aws_nat_eip_private_subnet2["tags"]["Name"]
  }
}

resource "aws_nat_gateway" "aws_nat_private_subnet2" {
  allocation_id = aws_eip.aws_nat_eip_private_subnet2.id
  subnet_id     = aws_subnet.aws_public_subnet2.id
  tags = {
    Name = var.aws_nat_private_subnet2["tags"]["Name"]
  }
  depends_on = [aws_internet_gateway.aws_internet_gateway]
}

resource "aws_route_table" "aws_route_table_private_subnet2" {
  vpc_id = aws_vpc.aws_vpc.id
  tags = {
    Name = var.aws_route_table_private_subnet2["tags"]["Name"]
  }
  route {
    cidr_block     = var.aws_route_table_private_subnet2["route"]["cidr_block"]
    nat_gateway_id = aws_nat_gateway.aws_nat_private_subnet2.id
  }
}

resource "aws_route_table_association" "aws_route_table_association_private_subnet2" {
  subnet_id      = aws_subnet.aws_private_subnet2.id
  route_table_id = aws_route_table.aws_route_table_private_subnet2.id
}


resource "aws_subnet" "aws_private_subnet3" {
  vpc_id            = aws_vpc.aws_vpc.id
  availability_zone = var.aws_private_subnet3["availability_zone"]
  cidr_block        = var.aws_private_subnet3["cidr_block"]
  tags = {
    Name = var.aws_private_subnet3["tags"]["Name"]
  }
}

resource "aws_eip" "aws_nat_eip_private_subnet3" {
  domain = var.aws_nat_eip_private_subnet3["domain"]
  tags = {
    Name = var.aws_nat_eip_private_subnet3["tags"]["Name"]
  }
}

resource "aws_nat_gateway" "aws_nat_private_subnet3" {
  allocation_id = aws_eip.aws_nat_eip_private_subnet3.id
  subnet_id     = aws_subnet.aws_public_subnet3.id
  tags = {
    Name = var.aws_nat_private_subnet3["tags"]["Name"]
  }
  depends_on = [aws_internet_gateway.aws_internet_gateway]
}

resource "aws_route_table" "aws_route_table_private_subnet3" {
  vpc_id = aws_vpc.aws_vpc.id
  tags = {
    Name = var.aws_route_table_private_subnet3["tags"]["Name"]
  }
  route {
    cidr_block     = var.aws_route_table_private_subnet3["route"]["cidr_block"]
    nat_gateway_id = aws_nat_gateway.aws_nat_private_subnet3.id
  }
}

## change nat gate way to subnet 3 once you have the eip

resource "aws_route_table_association" "aws_route_table_association_private_subnet3" {
  subnet_id      = aws_subnet.aws_private_subnet3.id
  route_table_id = aws_route_table.aws_route_table_private_subnet3.id
}
