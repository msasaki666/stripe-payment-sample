export default function AccountPage() {
  const portalUrl = new URL("/create-customer-portal-session", process.env.NEXT_PUBLIC_SERVER_BASE_URL).toString();

  return (
    <div>
      <h1>Account</h1>

      <form action={portalUrl} method="post">
        <button type="submit">決済ポータルへ</button>
      </form>
    </div>
  );
}
