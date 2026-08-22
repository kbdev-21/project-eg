export type User = {
  id: string;
  name: string;
  email: string | null;
  caroRating: number;
  role: string;
  avtCode: string;
  createdAt: Date;
  updatedAt: Date;
}