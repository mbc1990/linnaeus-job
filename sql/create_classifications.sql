CREATE TABLE classifications(
  class_id varchar(1024) primary key not null,
  class_name varchar(1024)  not null,
  image_id integer not null,
  probability decimal not null,
  timestamp timestamp default current_timestamp
);
