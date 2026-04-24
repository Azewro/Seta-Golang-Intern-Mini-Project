import { Link } from "react-router-dom";

export default function NotFoundPage() {
  return (
    <section className="card">
      <h2>Page not found</h2>
      <p>
        Go back to <Link to="/">home</Link>
      </p>
    </section>
  );
}

