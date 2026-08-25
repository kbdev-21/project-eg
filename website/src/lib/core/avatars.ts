import type { User } from "$lib/api/type";
import Bunny from "$lib/assets/Bunny.png";
import Kitty from "$lib/assets/Kitty.png";
import Teddy from "$lib/assets/Teddy.png";
import Hamster from "$lib/assets/Hamster.png";
import Monkey from "$lib/assets/Monkey.png";
import Piggy from "$lib/assets/Piggy.png";

export const AVATARS: Record<User["avtCode"], string> = {
  BUNNY: Bunny,
  KITTY: Kitty,
  TEDDY: Teddy,
  HAMSTER: Hamster,
  MONKEY: Monkey,
  PIGGY: Piggy,
};
