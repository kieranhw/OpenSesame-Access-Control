import { NextResponse } from "next/server";

export async function POST(request: Request) {
  let body: unknown;
  try {
    body = await request.json();
  } catch (err) {
    console.error("Invalid JSON in request:", err);
    return NextResponse.json(
      { error: "Invalid request body" },
      { status: 400 },
    );
  }

  let hubRes: Response;
  try {
    hubRes = await fetch("http://localhost:11072/management/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

    if (!hubRes.ok) {
      const text = await hubRes.text();
      return NextResponse.json(
        { error: "Login failed: " + text },
        { status: hubRes.status },
      );
    }
  } catch (err) {
    console.error("Unable to connect to the OpenSesame hub:", err);
    return NextResponse.json(
      { error: "Unable to connect to the OpenSesame hub" },
      { status: 502 },
    );
  }

  let hubData;
  try {
    hubData = await hubRes.json();
  } catch (err) {
    console.error("Invalid JSON from OpenSesame hub:", err);
    return NextResponse.json(
      { error: "Invalid response from the OpenSesame hub" },
      { status: 502 },
    );
  }

  const response = NextResponse.json(hubData, {
    status: hubRes.status,
  });

  const setCookie = hubRes.headers.get("set-cookie");
  if (setCookie) response.headers.set("Set-Cookie", setCookie);
  return response;
}
