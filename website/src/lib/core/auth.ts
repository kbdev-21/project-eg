import { createClient } from "@supabase/supabase-js";

const sbUrl = "https://meulllksjzrlelsjixzi.supabase.co";
const sbPublishableKey = "sb_publishable_-2N43aDwIHItePMeYybY3Q_87_bHNZc"

export const auth = createClient(sbUrl, sbPublishableKey, {
  auth: {
    flowType: "pkce"
  }
}).auth;