export type User = {
  id: string;
  name: string;
  normalizedName: string;
  email: string | null;
  caroRating: number;
  role: string;
  avtCode: "BUNNY" | "KITTY" | "TEDDY" | "HAMSTER" | "MONKEY" | "PIGGY";
  createdAt: Date;
  updatedAt: Date;
}