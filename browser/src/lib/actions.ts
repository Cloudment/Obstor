"use server";

import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { rpc } from "./rpc";

// Types
interface BucketSettings {
  name: string;
  versioning: boolean;
  objectLocking: boolean;
  accessPolicy: "private" | "public-read" | "public-read-write" | "custom";
  customPolicy: string;
  sftpEnabled: boolean;
  s3Enabled: boolean;
  anonymousAccess: boolean;
  quotaEnabled: boolean;
  quotaType: "hard" | "fifo";
  quotaSize: string;
  quotaUnit: "GB" | "TB" | "PB";
  encryptionEnabled: boolean;
  encryptionType: "SSE-S3" | "SSE-KMS";
  kmsKeyId: string;
  tags: { key: string; value: string }[];
}

// Auth
export async function loginAction(formData: FormData) {
  const accessKey = formData.get("accessKey") as string;
  const secretKey = formData.get("secretKey") as string;

  try {
    const result = await rpcUnauthed<{ token: string }>("Login", {
      username: accessKey,
      password: secretKey,
    });

    const cookieStore = await cookies();
    cookieStore.set("obstor_token", result.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
      maxAge: 60 * 60 * 24 * 7,
    });
  } catch {
    return { error: "Invalid access key or secret key" };
  }

  redirect("/");
}

async function rpcUnauthed<T>(method: string, params: Record<string, unknown> = {}): Promise<T> {
  const endpoint = process.env.OBSTOR_ENDPOINT || "http://localhost:9000";
  const res = await fetch(`${endpoint}/obstor/webrpc`, {
    method: "POST",
    headers: { "Content-Type": "application/json", "User-Agent": "Mozilla/5.0 Obstor Dashboard" },
    body: JSON.stringify({ id: 1, jsonrpc: "2.0", method: `web.${method}`, params }),
    cache: "no-store",
  });
  const data = await res.json();
  if (data.error) throw new Error(data.error.message);
  return data.result as T;
}

export async function logoutAction() {
  const cookieStore = await cookies();
  cookieStore.delete("obstor_token");
  redirect("/login");
}

export async function changePasswordAction(formData: FormData) {
  const currentAccessKey = formData.get("currentAccessKey") as string;
  const currentSecretKey = formData.get("currentSecretKey") as string;
  const newAccessKey = formData.get("newAccessKey") as string;
  const newSecretKey = formData.get("newSecretKey") as string;

  try {
    await rpc("SetAuth", { currentAccessKey, currentSecretKey, newAccessKey, newSecretKey });
    return { success: true };
  } catch (err) {
    return { error: err instanceof Error ? err.message : "Failed to change password" };
  }
}

// Buckets CRUD
export async function createBucketAction(formData: FormData) {
  const bucketName = formData.get("bucketName") as string;
  try {
    await rpc("MakeBucket", { bucketName });
    return { success: true, bucketName };
  } catch (err) {
    return { error: err instanceof Error ? err.message : "Failed to create bucket" };
  }
}

export async function deleteBucketAction(bucketName: string) {
  try {
    await rpc("DeleteBucket", { bucketName });
  } catch (err) {
    return { error: err instanceof Error ? err.message : "Failed to delete bucket" };
  }
  redirect("/");
}

export async function setBucketPolicyAction(formData: FormData) {
  const bucketName = formData.get("bucketName") as string;
  const prefix = formData.get("prefix") as string;
  const policy = formData.get("policy") as string;

  try {
    await rpc("SetBucketPolicy", { bucketName, prefix, policy });
    return { success: true };
  } catch (err) {
    return { error: err instanceof Error ? err.message : "Failed to set policy" };
  }
}

// Buckets Settings
function accessPolicyToBackend(
  policy: "private" | "public-read" | "public-read-write" | "custom",
): string {
  switch (policy) {
    case "public-read":
      return "readonly";
    case "public-read-write":
      return "readwrite";
    case "custom":
      return "readwrite";
    default:
      return "none";
  }
}

function backendPolicyToUI(
  policy: string,
): "private" | "public-read" | "public-read-write" | "custom" {
  switch (policy) {
    case "readonly":
      return "public-read";
    case "readwrite":
      return "public-read-write";
    case "writeonly":
      return "public-read-write";
    default:
      return "private";
  }
}

export async function getBucketSettingsAction(
  bucketName: string,
): Promise<BucketSettings | { error: string }> {
  try {
    let policyValue = "none";
    try {
      const res = await rpc<{ policy: string }>("GetBucketPolicy", {
        bucketName,
        prefix: "",
      });
      policyValue = res.policy || "none";
    } catch {
      // bucket may not have a policy set
    }

    const accessPolicy = backendPolicyToUI(policyValue);

    return {
      name: bucketName,
      versioning: false,
      objectLocking: false,
      accessPolicy,
      customPolicy: "",
      sftpEnabled: false,
      s3Enabled: true,
      anonymousAccess: accessPolicy !== "private",
      quotaEnabled: false,
      quotaType: "hard",
      quotaSize: "",
      quotaUnit: "GB",
      encryptionEnabled: false,
      encryptionType: "SSE-S3",
      kmsKeyId: "",
      tags: [],
    };
  } catch (err) {
    return {
      error: err instanceof Error ? err.message : "Failed to load bucket settings",
    };
  }
}

export async function createBucketWithSettingsAction(
  settings: BucketSettings,
): Promise<{ success: true; bucketName: string } | { error: string }> {
  try {
    await rpc("MakeBucket", { bucketName: settings.name });
    await applyBucketSettings(settings);
    revalidatePath("/");
    revalidatePath(`/${settings.name}`);
    return { success: true, bucketName: settings.name };
  } catch (err) {
    return {
      error: err instanceof Error ? err.message : "Failed to create bucket",
    };
  }
}

export async function updateBucketSettingsAction(
  settings: BucketSettings,
): Promise<{ success: true } | { error: string }> {
  try {
    await applyBucketSettings(settings);
    revalidatePath("/");
    revalidatePath(`/${settings.name}`);
    return { success: true };
  } catch (err) {
    return {
      error: err instanceof Error ? err.message : "Failed to update bucket settings",
    };
  }
}

async function applyBucketSettings(settings: BucketSettings) {
  const policyType = accessPolicyToBackend(settings.accessPolicy);
  await rpc("SetBucketPolicy", {
    bucketName: settings.name,
    prefix: "",
    policy: policyType,
  });
}

// Objects
export async function deleteObjectAction(bucketName: string, objectName: string) {
  try {
    await rpc("RemoveObject", { bucketname: bucketName, objects: [objectName] });
    return { success: true };
  } catch (err) {
    return { error: err instanceof Error ? err.message : "Failed to delete object" };
  }
}

export async function getShareLink(bucketName: string, objectName: string, expiry = 86400) {
  try {
    const result = await rpc<{ url: string }>("PresignedGet", {
      host: process.env.OBSTOR_HOST || "localhost:9000",
      bucket: bucketName,
      object: objectName,
      expiry,
    });
    return { url: result.url };
  } catch (err) {
    return { error: err instanceof Error ? err.message : "Failed to generate link" };
  }
}

export async function getUploadURL(bucketName: string, prefix: string, objectName: string) {
  try {
    const result = await rpc<{ url: string }>("PresignedPut", {
      host: process.env.OBSTOR_HOST || "localhost:9000",
      bucket: bucketName,
      prefix,
      object: objectName,
    });
    return { url: result.url };
  } catch (err) {
    return { error: err instanceof Error ? err.message : "Failed to get upload URL" };
  }
}

export async function getObjectChecksums(bucketName: string, objectName: string) {
  try {
    const result = await rpc<{ md5: string; sha1: string; sha256: string; sha512: string }>(
      "GetObjectChecksums",
      { bucketName, objectName },
    );
    return result;
  } catch (err) {
    return { error: err instanceof Error ? err.message : "Failed to compute checksums" };
  }
}

export async function getDownloadURL(bucketName: string, objectName: string) {
  try {
    const result = await rpc<{ url: string }>("PresignedGet", {
      host: process.env.OBSTOR_HOST || "localhost:9000",
      bucket: bucketName,
      object: objectName,
      expiry: 300,
    });
    return { url: result.url };
  } catch (err) {
    return { error: err instanceof Error ? err.message : "Failed to get download URL" };
  }
}
