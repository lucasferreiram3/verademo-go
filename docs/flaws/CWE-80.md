# CWE-80: Improper Neutralization of Script-Related HTML Tags in a Web Page (Basic XSS)
Verademo GO makes a call to a net:http::Error(), which contains a cross-site scripting flaw. The application populates the response with input that is untrusted and not verified. This as a result, allows attackers to embed any content they want such as JavaScript code, which can be executed in the victim's browser. XSS vulnerabilities are commonly used to steal or manipulate cookies, modify content, and compromise confidential information.

# Remediate
* Contextual escaping on all untrusted data before using it to construct any HTTP response can prevent this. Can use either entity escaping or attribute escaping depending on what the end goal is.

# Resources 
* [CWE-80](https://cwe.mitre.org/data/definitions/80.html)