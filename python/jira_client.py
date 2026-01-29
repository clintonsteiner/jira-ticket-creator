"""
JIRA Ticket Creator Python Client

Direct integration with Go backend via CGO C bindings.
No subprocess or REST API required - direct function calls.
"""

import ctypes
import json
import os
from pathlib import Path
from typing import Dict, Optional, List


class JiraClient:
    """Client for JIRA Ticket Creator using direct C bindings"""

    def __init__(
        self,
        url: str,
        email: str,
        token: str,
        project: Optional[str] = None,
        ticket: Optional[str] = None,
        lib_path: Optional[str] = None
    ):
        """Initialize JIRA client"""
        self.url = url
        self.email = email
        self.token = token

        # Resolve project from either project key or ticket key
        if project:
            self.project = project
        elif ticket:
            self.project = self._extract_project_sync(ticket)
        else:
            raise ValueError("Either 'project' or 'ticket' parameter is required")

        # Load the C library
        self.lib = self._load_library(lib_path)
        self._setup_function_signatures()

    def _load_library(self, lib_path: Optional[str]) -> ctypes.CDLL:
        """Load the compiled C library"""
        if lib_path and Path(lib_path).exists():
            return ctypes.CDLL(lib_path)

        # Try common locations
        candidates = [
            Path(__file__).parent / "libjira.so",
            Path(__file__).parent / "libjira.dylib",
            Path(__file__).parent / "libjira.dll",
            Path("/usr/local/lib/libjira.so"),
            Path("/usr/lib/libjira.so"),
        ]

        for candidate in candidates:
            if candidate.exists():
                return ctypes.CDLL(str(candidate))

        raise FileNotFoundError("Could not find libjira library. Run: make python-build")

    def _setup_function_signatures(self):
        """Setup ctypes function signatures"""
        self.lib.CreateTicketJSON.argtypes = [
            ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p,
            ctypes.c_char_p, ctypes.c_char_p
        ]
        self.lib.CreateTicketJSON.restype = ctypes.c_char_p

        self.lib.ExtractProjectKey.argtypes = [ctypes.c_char_p]
        self.lib.ExtractProjectKey.restype = ctypes.c_char_p

        self.lib.FreeMemory.argtypes = [ctypes.c_char_p]
        self.lib.FreeMemory.restype = None

        self.lib.Version.argtypes = []
        self.lib.Version.restype = ctypes.c_char_p

    def _extract_project_sync(self, ticket_key: str) -> str:
        """Extract project key from ticket key"""
        result = self.lib.ExtractProjectKey(ticket_key.encode())
        project = result.decode()
        self.lib.FreeMemory(result)

        if not project:
            raise ValueError(f"Invalid ticket key: {ticket_key}")

        return project

    def create_ticket(
        self,
        summary: str,
        description: str = "",
        issue_type: str = "Task",
        priority: str = "Medium",
        assignee: str = "",
        labels: Optional[List[str]] = None,
        blocked_by: Optional[List[str]] = None
    ) -> Dict[str, str]:
        """Create a JIRA ticket"""
        if not summary:
            raise ValueError("summary is required")

        request = {
            "summary": summary,
            "description": description,
            "issue_type": issue_type,
            "priority": priority,
        }

        if assignee:
            request["assignee"] = assignee
        if labels:
            request["labels"] = labels
        if blocked_by:
            request["blocked_by"] = blocked_by

        json_str = json.dumps(request)

        result = self.lib.CreateTicketJSON(
            self.url.encode(),
            self.email.encode(),
            self.token.encode(),
            self.project.encode(),
            json_str.encode()
        )

        response_json = result.decode()
        self.lib.FreeMemory(result)

        response = json.loads(response_json)

        if "error" in response and response["error"]:
            raise Exception(f"Failed to create ticket: {response['error']}")

        return response

    def extract_project_key(self, ticket_key: str) -> str:
        """Extract project key from ticket key"""
        return self._extract_project_sync(ticket_key)

    def get_version(self) -> str:
        """Get library version"""
        result = self.lib.Version()
        version = result.decode()
        self.lib.FreeMemory(result)
        return version


if __name__ == "__main__":
    try:
        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="api-token",
            project="PROJ"
        )
        print(f"Library version: {client.get_version()}")
    except Exception as e:
        print(f"Error: {e}")
