'use client';
import useSWR from "swr";

export default function ProductsPage() {
  const productsUrl = new URL("/products", process.env.NEXT_PUBLIC_SERVER_BASE_URL).toString();
  const { data, error, isLoading } = useSWR(
    productsUrl,
    (url) => fetch(url).then(res => res.json()).catch(err => console.error(err)),
  );

  if (error) return <div>error</div>
  if (isLoading) return <div>loading...</div>
  return (
    <div>
      <h1>Products</h1>
      <ul>
        {data.map((product: any) => (
          <li key={product.id}>
            <form action={new URL("/products/create-checkout-session", process.env.NEXT_PUBLIC_SERVER_BASE_URL).toString()} method="post">
              <p>{product.name}</p>
              <input type="hidden" name="stripe_price_id" value={product.stripe_price_id} />
              <button type="submit">購入する</button>
            </form>
          </li>
        ))}
      </ul>
    </div>
  );
}
