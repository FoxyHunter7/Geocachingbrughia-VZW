<?php

namespace Database\Seeders;

use App\Models\MessageLanguage;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class MessageLanguageSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        DB::table('message_languages')->delete();

        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/messagesLanguage.csv'), ';');

        $model = new MessageLanguage();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
