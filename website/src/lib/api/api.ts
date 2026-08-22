import { PUBLIC_API_BASE_URL } from "$env/static/public";
import axios from "axios";
import type { User } from "./type";

const baseUrl = PUBLIC_API_BASE_URL;

export async function getMe(token: string) {
  const res = await axios.get<User>(`${baseUrl}/api/me`, {
    headers: {
      Authorization: "Bearer " + token,
    },
  });
  return res.data;
}