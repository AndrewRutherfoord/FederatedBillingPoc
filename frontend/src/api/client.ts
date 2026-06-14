import createClient from "openapi-fetch";
import type { paths } from "@/api/api-schema";

const client = createClient<paths>({ baseUrl: "http://localhost:8082" });

export default client;