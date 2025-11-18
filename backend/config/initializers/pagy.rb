# Pagy initializer file
# See https://ddnexus.github.io/pagy/docs/api/pagy/

# Pagy default: 20 items per page
Pagy::DEFAULT[:limit] = 12

# Set the default page parameter name
Pagy::DEFAULT[:page_param] = :page

# Set up metadata keys for API responses
Pagy::DEFAULT[:metadata] = [
  :page,          # current page number
  :limit,         # items per page
  :count,         # total items count
  :pages,         # total pages count
  :prev_page,     # previous page number or nil
  :next_page      # next page number or nil
]
