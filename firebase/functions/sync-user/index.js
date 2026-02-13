// =============================================================================
// Firebase Cloud Function: sync-user
// =============================================================================
// Triggered on Firebase Auth user creation. Syncs the new user record into
// the SONA platform's PostgreSQL database via the Go backend REST API.
//
// Deploy with:
//   firebase deploy --only functions
//
// Environment variables (set via Firebase config):
//   sona.api_base_url  - Base URL of the Go backend (e.g. https://api.machanirobotics.dev)
// =============================================================================

const functions = require("firebase-functions");
const admin = require("firebase-admin");
const https = require("https");
const http = require("http");
const { URL } = require("url");

// Initialize Firebase Admin SDK (uses default service account in Cloud Functions)
admin.initializeApp();

/**
 * Sends a POST request to the Go backend API to create / sync a user record.
 *
 * @param {string} apiUrl - Full URL including path (e.g. https://api.example.com/api/auth/sync-user)
 * @param {object} payload - JSON body to send
 * @returns {Promise<object>} Parsed JSON response
 */
function postJSON(apiUrl, payload) {
  return new Promise((resolve, reject) => {
    const parsed = new URL(apiUrl);
    const transport = parsed.protocol === "https:" ? https : http;

    const body = JSON.stringify(payload);

    const options = {
      hostname: parsed.hostname,
      port: parsed.port || (parsed.protocol === "https:" ? 443 : 80),
      path: parsed.pathname + parsed.search,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Content-Length": Buffer.byteLength(body),
      },
    };

    const req = transport.request(options, (res) => {
      let data = "";
      res.on("data", (chunk) => {
        data += chunk;
      });
      res.on("end", () => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          try {
            resolve(JSON.parse(data));
          } catch (_) {
            resolve({ raw: data });
          }
        } else {
          reject(
            new Error(
              `API responded with status ${res.statusCode}: ${data}`
            )
          );
        }
      });
    });

    req.on("error", (err) => reject(err));
    req.write(body);
    req.end();
  });
}

/**
 * Cloud Function: onUserCreated
 *
 * Fires whenever a new user is created in Firebase Authentication.
 * Extracts the uid, email, and displayName and POSTs them to the Go backend
 * so the user is persisted in PostgreSQL.
 */
exports.onUserCreated = functions.auth.user().onCreate(async (user) => {
  const { uid, email, displayName } = user;

  functions.logger.info("New Firebase Auth user created", {
    uid,
    email,
    displayName,
  });

  // Resolve the backend API base URL from Firebase environment config.
  // Set it with:  firebase functions:config:set sona.api_base_url="https://api.machanirobotics.dev"
  const apiBaseUrl =
    (functions.config().sona && functions.config().sona.api_base_url) ||
    "http://api:3000";

  const syncEndpoint = `${apiBaseUrl}/api/auth/sync-user`;

  const payload = {
    firebase_uid: uid,
    email: email || null,
    display_name: displayName || null,
  };

  try {
    const result = await postJSON(syncEndpoint, payload);
    functions.logger.info("User synced to backend successfully", {
      uid,
      response: result,
    });
  } catch (error) {
    // Log the error but do NOT rethrow -- rethrowing would cause Cloud
    // Functions to retry the invocation, which may not be desirable for
    // user-creation events.  Adjust retry behaviour as needed.
    functions.logger.error("Failed to sync user to backend", {
      uid,
      error: error.message,
    });
  }
});
