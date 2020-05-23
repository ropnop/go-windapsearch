#!/usr/bin/env python

import requests
import sys
import json
from bs4 import BeautifulSoup

WIN32_ATTRIBUTES_PAGE = "https://docs.microsoft.com/en-us/windows/win32/adschema/"

ATTRIBUTES = []
TOTAL = 0
COUNTER = 0


def run(filename):
    global TOTAL
    global COUNTER
    w32page = requests.get(WIN32_ATTRIBUTES_PAGE+"/attributes-all")
    soup = BeautifulSoup(w32page.content, 'html.parser')
    datalinks = soup.find_all('dl')[0]
    hrefs = datalinks.find_all('a')
    TOTAL = len(hrefs)
    print("[+] Found {} attributes...".format(TOTAL))

    for href in hrefs:
        attributelink = WIN32_ATTRIBUTES_PAGE + "/" + href.attrs['href']
        getAttributeInfo(attributelink)
    
    with open(filename, "w") as fp:
        fp.write(json.dumps(ATTRIBUTES))
    print("{} written".format(filename))



def getAttributeInfo(link):
    global TOTAL
    global COUNTER
    COUNTER += 1
    attr = dict()
    attrPage = requests.get(link)
    if attrPage.status_code != 200:
        print("\t[!] Error fetching page: {}".format(link))
        return
    soup = BeautifulSoup(attrPage.content, 'html.parser')

    tag = soup.find("td", string="Ldap-Display-Name")
    if tag is None:
        print("\t[!] Didn't find 'Ldap-Display-Name' on {}".format(link))
        return
    
    infoTable = tag.find_parents('table')
    if len(infoTable) != 1:
        print("\t[!] Didn't find table on {}".format(link))
        return
    
    for row in infoTable[0].find_all('tr'):
        tds = row.find_all('td')
        if len(tds) != 2:
            continue
        name = tds[0].getText()
        value = tds[1].getText()
        attr[name] = value
    isSingleTags = soup.find_all("td", string="Is-Single-Valued")

    if len(isSingleTags) != 0:
        latest = isSingleTags[-1]
        row = latest.parent
        tds = row.find_all('td')
        if len(tds) == 2:
            attr["Is-Single-Valued"] = tds[1].getText() == 'True'
    import ipdb; ipdb.set_trace()
    ATTRIBUTES.append(attr)
    print("[+] ({}/{}) Fetched {}".format(COUNTER, TOTAL, attr.get('Ldap-Display-Name', "Unknown")))



if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("[+] Usage: {} output.json".format(sys.argv[0]))
        sys.exit(1)
    filename = sys.argv[1]
    run(filename)