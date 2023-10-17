files = ["./0030-31.json", 
"./0031-01.json",
 "1-2.json",
  "2-3.json",
   "3-4.json",
    "4-5.json",
     "5-6.json",
      "6-7.json",
       "7-8.json",
        "8-9.json",
         "9-10.json",
          "10-11.json",
           "11-12.json",
            "12-13.json",
             "13-14.json",
              "14-15.json",
               "15-16.json",
                "16-17.json",
                 "17-18.json",
                  "18-19.json",
                   "19-20.json",
                    "20-21.json",
                     "21-22.json",
                      "22-23.json",
                       "23-24.json",
                        "24-25.json",
                         "25-26.json",
                          "26-27.json",
                           "27-28.json"
 ]
exit_filename = "merged_json/merged.json"

import json


def merge_json_files(file_paths):
    merged_contents = []

    for file_path in file_paths:
        with open(file_path, 'r', encoding='utf-8') as file_in:
            merged_contents.extend(json.load(file_in))

    with open(exit_filename, 'w', encoding='utf-8') as file_out:
        json.dump(merged_contents, file_out)


paths = [
    'employees_1.json',
    'employees_2.json'
]

merge_json_files(files) 
