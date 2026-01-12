# Step 1: Use Node.js as the base image
FROM node:18-alpine AS build

# Step 2: Set the working directory
WORKDIR /app

# Step 3: Copy package.json and lock file
COPY package.json pnpm-lock.yaml ./

# Step 4: Install dependencies
RUN npm install -g pnpm && pnpm install --frozen-lockfile

# Step 5: Copy the rest of the application files
COPY . .

# Step 6: Build the frontend application
RUN pnpm build

# Step 7: Use a lightweight web server for serving the app
FROM nginx:1.25-alpine

# Step 8: Copy the built application from the previous stage
COPY --from=build /app/.output/public /usr/share/nginx/html

# Step 9: Expose port 3000 and start Nginx
EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]
