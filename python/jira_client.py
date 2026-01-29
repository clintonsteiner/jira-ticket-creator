"""
JIRA Ticket Creator Python Client

Direct integration with Go backend via CGO C bindings.
No subprocess or REST API required - direct function calls.
"""

import ctypes
import json
import os
from pathlib import Path
from typing import Dict, Optional, List, Any


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
        # CreateTicket function: (url, email, token, project, json) -> json_response
        self.lib.CreateTicketJSON.argtypes = [
            ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p,
            ctypes.c_char_p, ctypes.c_char_p
        ]
        self.lib.CreateTicketJSON.restype = ctypes.c_char_p

        # GetTicket function: (url, email, token, key) -> json_response
        self.lib.GetTicket.argtypes = [
            ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p
        ]
        self.lib.GetTicket.restype = ctypes.c_char_p

        # SearchTickets function: (url, email, token, jql) -> json_response
        self.lib.SearchTickets.argtypes = [
            ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p
        ]
        self.lib.SearchTickets.restype = ctypes.c_char_p

        # UpdateTicket function: (url, email, token, key, json) -> json_response
        self.lib.UpdateTicket.argtypes = [
            ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p,
            ctypes.c_char_p, ctypes.c_char_p
        ]
        self.lib.UpdateTicket.restype = ctypes.c_char_p

        # ExtractProjectKey function: (ticket_key) -> project_key
        self.lib.ExtractProjectKey.argtypes = [ctypes.c_char_p]
        self.lib.ExtractProjectKey.restype = ctypes.c_char_p

        # FreeMemory function: (ptr) -> void
        self.lib.FreeMemory.argtypes = [ctypes.c_char_p]
        self.lib.FreeMemory.restype = None

        # Version function: () -> version_string
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

    def get_ticket(self, ticket_key: str) -> Dict[str, Any]:
        """Get ticket details by key

        Args:
            ticket_key: The ticket key (e.g., 'PROJ-123')

        Returns:
            Dictionary with ticket details: key, id, summary, description, status,
            issue_type, priority, assignee, labels, url
        """
        result = self.lib.GetTicket(
            self.url.encode(),
            self.email.encode(),
            self.token.encode(),
            ticket_key.encode()
        )

        response_json = result.decode()
        self.lib.FreeMemory(result)
        response = json.loads(response_json)

        if "error" in response and response["error"]:
            raise Exception(f"Failed to get ticket: {response['error']}")

        return response

    def search(self, jql: str = "", **kwargs) -> Dict[str, Any]:
        """Search for tickets using JQL or keyword arguments

        Args:
            jql: JQL query string (e.g., 'project = PROJ AND status = "In Progress"')
            **kwargs: Alternative keyword arguments:
                - key: Search by ticket key
                - summary: Search by summary text
                - status: Filter by status
                - assignee: Filter by assignee
                - issue_type: Filter by issue type

        Returns:
            Dictionary with search results: total, count, tickets
        """
        # Build JQL if kwargs provided
        if not jql and kwargs:
            jql_parts = []
            if "key" in kwargs:
                jql_parts.append(f'key = {kwargs["key"]}')
            if "summary" in kwargs:
                jql_parts.append(f'summary ~ "{kwargs["summary"]}"')
            if "status" in kwargs:
                jql_parts.append(f'status = "{kwargs["status"]}"')
            if "assignee" in kwargs:
                jql_parts.append(f'assignee = "{kwargs["assignee"]}"')
            if "issue_type" in kwargs:
                jql_parts.append(f'type = {kwargs["issue_type"]}')

            if jql_parts:
                jql = " AND ".join(jql_parts)

        if not jql:
            jql = f"project = {self.project}"

        result = self.lib.SearchTickets(
            self.url.encode(),
            self.email.encode(),
            self.token.encode(),
            jql.encode()
        )

        response_json = result.decode()
        self.lib.FreeMemory(result)
        response = json.loads(response_json)

        if "error" in response and response["error"]:
            raise Exception(f"Search failed: {response['error']}")

        return response

    def update_ticket(self, ticket_key: str, **kwargs) -> Dict[str, Any]:
        """Update ticket fields

        Args:
            ticket_key: The ticket key to update
            **kwargs: Fields to update:
                - summary: New summary
                - description: New description
                - priority: New priority (Lowest, Low, Medium, High, Highest)
                - assignee: New assignee email

        Returns:
            Dictionary with update response
        """
        request = {}

        if "summary" in kwargs:
            request["summary"] = kwargs["summary"]
        if "description" in kwargs:
            request["description"] = kwargs["description"]
        if "priority" in kwargs:
            request["priority"] = kwargs["priority"]
        if "assignee" in kwargs:
            request["assignee"] = kwargs["assignee"]

        if not request:
            raise ValueError("No fields specified for update")

        json_str = json.dumps(request)

        result = self.lib.UpdateTicket(
            self.url.encode(),
            self.email.encode(),
            self.token.encode(),
            ticket_key.encode(),
            json_str.encode()
        )

        response_json = result.decode()
        self.lib.FreeMemory(result)
        response = json.loads(response_json)

        if "error" in response and response["error"]:
            raise Exception(f"Failed to update ticket: {response['error']}")

        return response

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
