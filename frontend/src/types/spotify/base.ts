export type ID = string;
export type Numeric = number;
export type URI = string;

export type Image = {
  height: Numeric;
  width: Numeric;
  url: string;
};

export type Followers = {
  total?: number | null;
  href: any;
};

export type Error = {
  status: number;
  message: string;
};
