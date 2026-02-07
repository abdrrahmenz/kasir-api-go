import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'https://laravel-kasir-go.3rkvle.easypanel.host';

export const options = {
  stages: [
    { duration: '30s', target: 20 },   // Ramp up to 20 VUs over 30s
    { duration: '1m30s', target: 50 }, // Hold at 50 VUs for 1m30s
    { duration: '30s', target: 0 },    // Ramp down to 0 VUs over 30s
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'], // 95th percentile should be <500ms, 99th <1000ms
    http_req_failed: ['rate<0.1'],                    // Error rate should be <10%
  },
};

export default function () {
  // Test GET all products
  let getProductsResponse = http.get(`${BASE_URL}/api/produk`);
  check(getProductsResponse, {
    'GET /api/produk - status is 200': (r) => r.status === 200,
    'GET /api/produk - response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1);

  // Extract product ID from response if available
  let productId = 1; // Default fallback
  if (getProductsResponse.status === 200) {
    try {
      const products = JSON.parse(getProductsResponse.body);
      if (products && products.data && products.data.length > 0) {
        productId = products.data[0].id;
      }
    } catch (e) {
      console.log('Failed to parse products response');
    }
  }

  // Test GET product by ID
  let getProductResponse = http.get(`${BASE_URL}/api/produk/${productId}`);
  check(getProductResponse, {
    'GET /api/produk/{id} - status is 200': (r) => r.status === 200,
    'GET /api/produk/{id} - response time < 300ms': (r) => r.timings.duration < 300,
  });

  sleep(1);

  // Test GET all categories
  let getCategoriesResponse = http.get(`${BASE_URL}/api/categories`);
  check(getCategoriesResponse, {
    'GET /api/categories - status is 200': (r) => r.status === 200,
    'GET /api/categories - response time < 300ms': (r) => r.timings.duration < 300,
  });

  sleep(2);
}
