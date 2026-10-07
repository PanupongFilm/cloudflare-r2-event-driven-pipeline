variable "cloudflare_api_token" {
  description = "Cloudflare API Token for managing R2, Queues, and Workers"
  type        = string
  sensitive   = true 
}

variable "cloudflare_account_id" {
  description = "Cloudflare Account ID"
  type        = string
}
